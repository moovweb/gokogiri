#include <libxml/tree.h>
#include "Callback.h"
#include <stdio.h>

void invalidTree(xmlNode *node, void *doc) {
    xmlNode *cur_node = NULL;
    if (node != NULL) {
        //fprintf(stderr, "invalid 0x%x\n", node);
        
        for (cur_node = node->children; cur_node; cur_node = cur_node->next) {
            //fprintf(stderr, "invalid child 0x%x\n", cur_node);
            invalidTree(cur_node, doc);
        }
    }
}
