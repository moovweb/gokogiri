#include "helper.h"

#include <time.h>
#include <sys/time.h>

#ifdef __MACH__
#include <mach/clock.h>
#include <mach/mach.h>
#endif

#include <pthread.h>
#include <stdio.h>
/* for ETIMEDOUT */
#include <errno.h>
#include <string.h>

void current_utc_time(struct timespec *ts) {
 
#ifdef __MACH__ // OS X does not have clock_gettime, use clock_get_time
  clock_serv_t cclock;
  mach_timespec_t mts;
  host_get_clock_service(mach_host_self(), CALENDAR_CLOCK, &cclock);
  clock_get_time(cclock, &mts);
  mach_port_deallocate(mach_task_self(), cclock);
  ts->tv_sec = mts.tv_sec;
  ts->tv_nsec = mts.tv_nsec;
#else
  clock_gettime(CLOCK_REALTIME, ts);
#endif
 
}

typedef struct xmlXpathEvalArgs_ {
    xmlXPathCompExprPtr comp;
    xmlXPathContextPtr ctx;
    xmlXPathObjectPtr result;
    pthread_cond_t done;
    int isDone;
} xmlXpathEvalArgs;


xmlNode* fetchNode(xmlNodeSet *nodeset, int index) {
    return nodeset->nodeTab[index];
}

xmlXPathObjectPtr go_resolve_variables(void* ctxt, char* name, char* ns);
int go_can_resolve_function(void* ctxt, char* name, char* ns);
void exec_xpath_function(xmlXPathParserContextPtr ctxt, int nargs);

xmlXPathFunction go_resolve_function(void* ctxt, char* name, char* ns) {
    // TODO(Noj) uncomment once this issue is resolved: https://code.google.com/p/go/issues/detail?id=6661
    //if (go_can_resolve_function(ctxt, name, ns))
    //    return exec_xpath_function;

    return 0;
}

void set_var_lookup(xmlXPathContext* c, void* data) {
    c->varLookupFunc = (void *)go_resolve_variables;
    c->varLookupData = data;
}

void set_function_lookup(xmlXPathContext* c, void* data) {
    c->funcLookupFunc = (void *)go_resolve_function;
    c->funcLookupData = data;
}

int getXPathObjectType(xmlXPathObject* o) {
    if(o == 0)
        return 0;
    return o->type;
}

void *expensiveCallXPathEval(void *data)
{
        int oldtype;

        /* allow the thread to be killed at any time */
        pthread_setcanceltype(PTHREAD_CANCEL_ASYNCHRONOUS, &oldtype);

        xmlXpathEvalArgs *args = (xmlXpathEvalArgs*) data;

        /* ... calculations and expensive io here, for example:
         * infinitely loop
         */
        args->result = xmlXPathCompiledEval(args->comp, args->ctx);
        args->isDone = 1;

        /* wake up the caller if we've completed in time */
        pthread_cond_signal(&args->done);
        return NULL;
}

xmlXPathObjectPtr xmlXPathEvalWithTimeout(xmlXPathCompExprPtr comp, xmlXPathContextPtr ctx, int timeout) {

    pthread_mutex_t calculating;
    pthread_cond_t done;
    pthread_mutex_init(&calculating, NULL);
    pthread_cond_init(&done, NULL);

    xmlXpathEvalArgs args;
    args.comp = comp;
    args.ctx = ctx;
    args.done = done;
    args.isDone = 0;
    args.result = NULL;

    struct timespec abs_time;
    pthread_t tid;
    int err;

    pthread_mutex_lock(&calculating);

    /* pthread cond_timedwait expects an absolute time to wait until */
    // clock_gettime(CLOCK_REALTIME, &abs_time);
    current_utc_time(&abs_time); 

    struct timespec max_wait;
    memset(&max_wait, 0, sizeof(max_wait));
    max_wait.tv_sec = timeout;

    abs_time.tv_sec += max_wait.tv_sec;
    abs_time.tv_nsec += max_wait.tv_nsec;

    pthread_create(&tid, NULL, expensiveCallXPathEval, &args);

    while (!args.isDone) {
        err = pthread_cond_timedwait(&done, &calculating, &abs_time);
    }

    if (err == ETIMEDOUT && !args.isDone) {
        fprintf(stderr, "%s: error and calculation timed out\n", __func__);
        pthread_mutex_unlock(&calculating);
        return NULL;
    } else {
        fprintf(stderr, "%s: no error\n", __func__);
        pthread_mutex_unlock(&calculating);
    }
            
    return args.result;

}