#include <stdio.h>
#include <string.h>
#include <libxml/parser.h>
#include <libxml/tree.h>
#include <libxml/xmlstring.h>

#include <cppunit/extensions/TestFactoryRegistry.h>
#include <cppunit/ui/text/TestRunner.h>
#include <cppunit/CompilerOutputter.h>
#include <cppunit/TestCase.h>
#include <cppunit/extensions/HelperMacros.h>

extern "C" {
#include "Chelpers.h"
}

#define XML_STR "<foo><bar>foo.bar</bar><div>div</div></foo>"

using namespace std;

class ChelpersTest : public CppUnit::TestCase  {
    CPPUNIT_TEST_SUITE( ChelpersTest );
    CPPUNIT_TEST( test_xml_output );
    CPPUNIT_TEST( test_xmlElement_append );
    CPPUNIT_TEST( test_xmlElement_prepend );
    CPPUNIT_TEST_SUITE_END();

    void test_xml_output();
    void test_xmlElement_append();
    void test_xmlElement_prepend();
    const char* next_line(const char *xml);
    private:
        xmlDocPtr doc;
        xmlChar* xmlbuff;
    public:
        void setUp();
        void tearDown();
};

void ChelpersTest::setUp() {
    xmlInitParser ();
    char *content = (char *)XML_STR;
    doc = xmlReadMemory(content, strlen(content), "", "UTF-8", 0);
}

void ChelpersTest::tearDown() {
    //printf("allocated: %d\n", xmlMemBlocks());
    //xmlMemDisplay(stdout);
    //printf("\n\n");
    if (xmlbuff != NULL)
        xmlFree(xmlbuff);
    if (doc != NULL)
        xmlFreeDoc(doc);
    //xmlCleanupMemory();
    xmlCleanupParser ();
    //xmlMemoryDump();
    xmlMemDisplay(stdout);
    printf("allocated: %d\n", xmlMemBlocks());
}


const char* ChelpersTest::next_line(const char *xml) {
    char* chr = strchr((char *) xmlbuff, '\n');
    if (chr == NULL)
        return xml;
    chr ++;
    return (const char*) chr;
}

void ChelpersTest::test_xml_output() {
  int buffersize;
  xmlDocDumpFormatMemory(doc, &xmlbuff, &buffersize, 0);
  const char* next = next_line((char *) xmlbuff);
  CPPUNIT_ASSERT(strncmp(next, XML_STR, strlen(XML_STR)) == 0);

}

void ChelpersTest::test_xmlElement_append() {
    xmlNodePtr root_element = xmlDocGetRootElement(doc);
    string insert_content = "<moov>z</moov>hello<web>c</web>";
    xmlElement_append(root_element, doc, insert_content.c_str(), insert_content.length(), NULL);
    int buffersize;
    xmlDocDumpFormatMemory(doc, &xmlbuff, &buffersize, 0);
    const char* next = next_line((char *) xmlbuff);
    //printf("%s",  (char*) next);
    string expected_xml = "<foo><bar>foo.bar</bar><div>div</div><moov>z</moov>hello<web>c</web></foo>";
    CPPUNIT_ASSERT(strncmp(next, expected_xml.c_str(), expected_xml.length()) == 0);
}

void ChelpersTest::test_xmlElement_prepend() {
    xmlNodePtr root_element = xmlDocGetRootElement(doc);
    string insert_content = "<moov>z</moov>hello<web>c</web>";
    xmlElement_prepend(root_element, doc, insert_content.c_str(), insert_content.length(), NULL);
    int buffersize;
    xmlDocDumpFormatMemory(doc, &xmlbuff, &buffersize, 0);
    const char* next = next_line((char *) xmlbuff);
    //printf("%s",  (char*) next);
    string expected_xml = "<foo><moov>z</moov>hello<web>c</web><bar>foo.bar</bar><div>div</div></foo>";
    CPPUNIT_ASSERT(strncmp(next, expected_xml.c_str(), expected_xml.length()) == 0);

}

CPPUNIT_TEST_SUITE_NAMED_REGISTRATION( ChelpersTest, "ChelpersTest" );

CppUnit::Test *suite()  {
    CppUnit::TestFactoryRegistry &registry =
    CppUnit::TestFactoryRegistry::getRegistry();

    registry.registerFactory(
      &CppUnit::TestFactoryRegistry::getRegistry( "ChelpersTest" ) );
    return registry.makeTest();
}

int main( int argc, char* argv[] ) {
  // if command line contains "-selftest" then this is the post build check
  // => the output must be in the compiler error format.
  bool selfTest = (argc > 1)  &&
    (std::string("-selftest") == argv[1]);

  CppUnit::TextUi::TestRunner runner;
  runner.addTest( suite() );   // Add the top suite to the test runner

  if ( selfTest )
  { // Change the default outputter to a compiler error format outputter
    // The test runner owns the new outputter.
    runner.setOutputter( CppUnit::CompilerOutputter::defaultOutputter(
          &runner.result(),
          std::cerr ) );
  }

  // Run the test.
  bool wasSucessful = runner.run( "" );

  // Return error code 1 if any tests failed.
  return wasSucessful ? 0 : 1;
}
