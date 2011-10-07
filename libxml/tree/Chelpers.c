#include <stdio.h>
#include <stdlib.h>

#include <libxml/parser.h>
#include <libxml/tree.h>
#include "Chelpers.h"

void xmlNode_print(xmlNodePtr node) {
    if (node->type == XML_ELEMENT_NODE) {
        printf("node type: Element, name: %s\n", node->name);
    }
}

int xmlElement_append(xmlNodePtr node, xmlDocPtr doc, const char* content, int content_length, const char* encoding) {
    xmlDocPtr new_doc = NULL;
    xmlNodePtr root_element = NULL;
    xmlNodePtr cur_node = NULL;
    xmlNodePtr next_node = NULL;

    char *wrapped_content = NULL;
    int wrapped_content_length = content_length;

    if (content_length <= 0) 
        return 0;

    //wrapped the content with <root></root>
    wrapped_content_length += 13;
    wrapped_content = (char*) malloc((wrapped_content_length+1)*sizeof(char));
    snprintf(wrapped_content, wrapped_content_length+1, "<root>%s</root>", content);
    //printf("content:%s %d", content, content_length); 
    //printf("wrapped content: %s %d\n", wrapped_content, wrapped_content_length);

    new_doc = xmlReadMemory(wrapped_content, wrapped_content_length, "", encoding, 0);
    if (new_doc == NULL) {
        free(wrapped_content);
        return 0;
    }

    root_element = xmlDocGetRootElement(new_doc);
    for (cur_node = root_element->children; cur_node; cur_node = next_node) {
        next_node = cur_node->next;
        xmlUnlinkNode(cur_node);
        cur_node = xmlDocCopyNode(cur_node, doc, 1);
        xmlAddChild(node, cur_node);
    }

    if (new_doc != NULL)
        xmlFreeDoc(new_doc);
    free(wrapped_content);
    return 1;
}

int xmlElement_prepend(xmlNodePtr node, xmlDocPtr doc, const char* content, int content_length, const char* encoding) {
    xmlDocPtr new_doc = NULL;
    xmlNodePtr root_element = NULL;
    xmlNodePtr cur_node = NULL;
    xmlNodePtr next_node = NULL;
    xmlNodePtr first_child = NULL;

    char *wrapped_content = NULL;
    int wrapped_content_length = content_length;

    if (content_length <= 0) 
        return 0;

    //wrapped the content with <root></root>
    wrapped_content_length += 13;
    wrapped_content = (char*) malloc((wrapped_content_length+1)*sizeof(char));
    snprintf(wrapped_content, wrapped_content_length+1, "<root>%s</root>", content);
    //printf("content:%s %d\n", content, content_length); 
    //printf("wrapped content: %s\n", wrapped_content);

    new_doc = xmlReadMemory(wrapped_content, wrapped_content_length, "", encoding, 0);
    if (new_doc == NULL) {
        free(wrapped_content);
        return 0;
    }

    root_element = xmlDocGetRootElement(new_doc);
    first_child = node->children;
    
    if (first_child == NULL) {
        for (cur_node = root_element->children; cur_node; cur_node = next_node) {
            next_node = cur_node->next;
            xmlUnlinkNode(cur_node);
            cur_node = xmlDocCopyNode(cur_node, doc, 1);
            xmlAddChild(node, cur_node);
        }
    }
    else {
        for (cur_node = root_element->children; cur_node; cur_node = next_node) {
            next_node = cur_node->next;
            xmlUnlinkNode(cur_node);
            cur_node = xmlDocCopyNode(cur_node, doc, 1);
            xmlAddPrevSibling(first_child, cur_node);
        }
    }

    if (new_doc != NULL)
        xmlFreeDoc(new_doc);
    free(wrapped_content);
    return 1;
}

