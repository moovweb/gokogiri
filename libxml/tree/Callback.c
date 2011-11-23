#include <libxml/tree.h>
#include "Callback.h"
#include <stdio.h>

void invalidTree(xmlNode *node, void *doc) {
    xmlNode *cur_node = NULL;
    if (node != NULL) {
        cur_node = node->children;
        while(cur_node != NULL) {
            invalidTree(cur_node, doc);
            cur_node = node->children;            
        }
        invalidNode(node, doc);
        if (node->parent != NULL) {
            xmlUnlinkNode(node);
        }
        xmlFreeNode(node);
    }
}
