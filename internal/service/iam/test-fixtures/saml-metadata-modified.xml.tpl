<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" xmlns:ds="http://www.w3.org/2000/09/xmldsig#" entityID="${entity_id_modified}" validUntil="2025-09-02T18:27:19.710Z">
   <md:IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
      <md:KeyDescriptor use="signing">
         <ds:KeyInfo>
            <ds:X509Data>
               <ds:X509Certificate>MIIErDCCA5SgAwIBAgIOAU+PT8RBAAAAAHxJXEcwDQYJKoZIhvcNAQELBQAwgZAxKDAmBgNVBAMMH1NlbGZTaWduZWRDZXJ0XzAyU2VwMjAxNV8xODI2NTMxGDAWBgNVBAsMDzAwRDI0MDAwMDAwcEFvQTEXMBUGA1UECgwOU2FsZXNmb3JjZS5jb20xFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xCzAJBgNVBAgMAkNBMQwwCgYDVQQGEwNVU0EwHhcNMTUwOTAyMTgyNjUzWhcNMTcwOTAyMTIwMDAwWjCBkDEoMCYGA1UEAwwfU2VsZlNpZ25lZENlcnRfMDJTZXAyMDE1XzE4MjY1MzEYMBYGA1UECwwPMDBEMjQwMDAwMDBwQW9BMRcwFQYDVQQKDA5TYWxlc2ZvcmNlLmNvbTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzELMAkGA1UECAwCQ0ExDDAKBgNVBAYTA1VTQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJp/wTRr9n1IWJpkRTjNpep47OKJrD2E6rGbJ18TG2RxtIz+zCn2JwH2aP3TULh0r0hhcg/pecv51RRcG7O19DBBaTQ5+KuoICQyKZy07/yDXSiZontTwkEYs06ssTwTHUcRXbcwTKv16L7omt0MjIhTTGfvtLOYiPwyvKvzAHg4eNuAcli0duVM78UIBORtdmy9C9ZcMh8yRJo5aPBq85wsE3JXU58ytyZzCHTBLH+2xFQrjYnUSEW+FOEEpI7o33MVdFBvWWg1R17HkWzcve4C30lqOHqvxBzyESZ/N1mMlmSt8gPFyB+mUXY99StJDJpnytbY8DwSzMQUo/sOVB0CAwEAAaOCAQAwgf0wHQYDVR0OBBYEFByu1EQqRQS0bYQBKS9K5qwKi+6IMA8GA1UdEwEB/wQFMAMBAf8wgcoGA1UdIwSBwjCBv4AUHK7URCpFBLRthAEpL0rmrAqL7oihgZakgZMwgZAxKDAmBgNVBAMMH1NlbGZTaWduZWRDZXJ0XzAyU2VwMjAxNV8xODI2NTMxGDAWBgNVBAsMDzAwRDI0MDAwMDAwcEFvQTEXMBUGA1UECgwOU2FsZXNmb3JjZS5jb20xFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xCzAJBgNVBAgMAkNBMQwwCgYDVQQGEwNVU0GCDgFPj0/EQQAAAAB8SVxHMA0GCSqGSIb3DQEBCwUAA4IBAQA9O5o1tC71qJnkq+ABPo4A1aFKZVT/07GcBX4/wetcbYySL4Q2nR9pMgfPYYS1j+P2E3viPsQwPIWDUBwFkNsjjX5DSGEkLAioVGKRwJshRSCSynMcsVZbQkfBUiZXqhM0wzvoa/ALvGD+aSSb1m+x7lEpDYNwQKWaUW2VYcHWv9wjujMyy7dlj8E/jqM71mw7ThNl6k4+3RQ802dMa14txm8pkF0vZgfpV3tkqhBqtjBAicVCaveqr3r3iGqjvyilBgdY+0NR8szqzm7CD/Bkb22+/IgM/mXQuL9KHD/WADlSGmYKmG3SSahmcZxznYCnzcRNN9LVuXlz5cbljmBj</ds:X509Certificate>
            </ds:X509Data>
         </ds:KeyInfo>
      </md:KeyDescriptor>
      <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
      <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="${entity_id_modified}/idp/endpoint/HttpPost"/>
      <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="${entity_id_modified}/idp/endpoint/HttpRedirect"/>
   </md:IDPSSODescriptor>
</md:EntityDescriptor>
