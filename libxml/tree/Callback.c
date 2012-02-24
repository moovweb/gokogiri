#include <libxml/tree.h>
#include "Callback.h"
#include <stdio.h>
#include <assert.h>

void clearXmlNode(xmlNode *node, void *doc) {
	xmlAttr *attr = NULL;
	//clear all attributes
	attr = node->properties;
	while (attr != NULL) {
	    //just mark it as deleted in the map
	    //all properties will be deleted along with this node
	    invalidateNode(attr, doc);
	    attr = attr->next;
	}
	invalidateNode(node, doc);
	//printf("clear node 0x%x\n", (unsigned int)node);
	if (node->parent != NULL) {
	    xmlUnlinkNode(node);
	}
	xmlFreeNode(node);
}

void invalidateTree(xmlNode *node, void *doc) {
    if (node != NULL) {
	    xmlNode *cur_node = node, *node_to_clear = node;
		do {
			while(cur_node->children != NULL) {
				cur_node = cur_node->children;
			}
			while(cur_node->children == NULL && cur_node != node) {
				node_to_clear = cur_node;
				cur_node = cur_node->parent;
				clearXmlNode(node_to_clear, doc);
			}
		} while (cur_node->children != NULL);

		assert(node == cur_node);
		clearXmlNode(cur_node, doc);
    }
}
