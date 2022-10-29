package zxml

import "testing"

func TestFormat(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{``, ``},
		{`<?xml?>`, `<?xml?>`},
		{`<elem>aa</elem>`, `<elem>aa</elem>`},

		{`
<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<greeting>
		<empty attr1="A" attr2="B"></empty>
		<svID>Test EPP server</svID>
		<svDate>2022-10-29T13:35:45+02:00</svDate>
		<svcMenu>
			<version>1.0</version>
			<version>2.0</version>
			<lang>en</lang>
			<objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
			<objURI>urn:ietf:params:xml:ns:host-1.0</objURI>
			<objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
			<objURI>urn:ietf:params:xml:ns:registrar-info-1.0</objURI>
			<objURI>urn:ietf:params:xml:ns:rgp-1.0</objURI>
			<objURI>urn:ietf:params:xml:ns:secDNS-1.1</objURI>
			<svcExtension>
				<extURI>http://www.subreg.cz/epp/gransy-domain-0.1</extURI>
			</svcExtension>
			<svcExtension>
				<extURI>http://www.subreg.cz/epp/gransy-document-0.1</extURI>
			</svcExtension>
			<svcExtension>
				<extURI>http://www.subreg.cz/epp/gransy-contact-0.1</extURI>
			</svcExtension>
		</svcMenu>
	</greeting>
</epp>`, `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
    <greeting>
        <empty attr1="A" attr2="B"></empty>
        <svID>Test EPP server</svID>
        <svDate>2022-10-29T13:35:45+02:00</svDate>
        <svcMenu>
            <version>1.0</version>
            <version>2.0</version>
            <lang>en</lang>
            <objURI>urn:ietf:params:xml:ns:domain-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:host-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:contact-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:registrar-info-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:rgp-1.0</objURI>
            <objURI>urn:ietf:params:xml:ns:secDNS-1.1</objURI>
            <svcExtension>
                <extURI>http://www.subreg.cz/epp/gransy-domain-0.1</extURI>
            </svcExtension>
            <svcExtension>
                <extURI>http://www.subreg.cz/epp/gransy-document-0.1</extURI>
            </svcExtension>
            <svcExtension>
                <extURI>http://www.subreg.cz/epp/gransy-contact-0.1</extURI>
            </svcExtension>
        </svcMenu>
    </greeting>
</epp>`,
		},

		{
			`<h:table><h:tr><h:td>Apples</h:td><h:td>Bananas</h:td></h:tr></h:table>
            <f:table><f:name>African Coffee Table</f:name><f:width>80</f:width><f:length>120</f:length></f:table>`,
			`<h:table>
    <h:tr>
        <h:td>Apples</h:td>
        <h:td>Bananas</h:td>
    </h:tr>
</h:table>
<f:table>
    <f:name>African Coffee Table</f:name>
    <f:width>80</f:width>
    <f:length>120</f:length>
</f:table>`,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have, err := Format([]byte(tt.in), "", "    ")
			if err != nil {
				t.Fatal(err)
			}

			// Don't use ztest.Diff with ztest.DiffXML, because this uses the
			// Format(), and we'd just be testing ourselves.
			if string(have) != string(tt.want) {
				t.Errorf("\nhave:\n%s\nwant:%s", have, tt.want)
			}
		})
	}
}
