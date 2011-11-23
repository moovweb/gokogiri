#include <libxml/tree.h>
#include "Callback.h"
#include <stdio.h>

void invalidTree(xmlNode *node, void *doc) {
    xmlNode *cur_node = NULL;
    xmlAttr *cur_attr = NULL;
    if (node != NULL) {
        //clear all attributes
        cur_attr = node->properties;
        while (cur_attr != NULL) {
            //just mark it as deleted in the map
            //all properties will be deleted along with this node
            invalidNode(cur_attr, doc);
            cur_attr = cur_attr->next;
        }
        //recursively clear all children
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
