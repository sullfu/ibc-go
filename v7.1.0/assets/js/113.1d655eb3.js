(window.webpackJsonp=window.webpackJsonp||[]).push([[113],{672:function(e,t,o){"use strict";o.r(t);var l=o(1),a=Object(l.a)({},(function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[o("h1",{attrs:{id:"clientstate"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#clientstate"}},[e._v("#")]),e._v(" "),o("code",[e._v("ClientState")])]),e._v(" "),o("p",[e._v("The 09-localhost "),o("code",[e._v("ClientState")]),e._v(" maintains a single field used to track the latest sequence of the state machine i.e. the height of the blockchain.")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"dHlwZSBDbGllbnRTdGF0ZSBzdHJ1Y3QgewogIC8vIHRoZSBsYXRlc3QgaGVpZ2h0IG9mIHRoZSBibG9ja2NoYWluCiAgTGF0ZXN0SGVpZ2h0IGNsaWVudHR5cGVzLkhlaWdodAp9Cg=="}}),e._v(" "),o("p",[e._v("The 09-localhost "),o("code",[e._v("ClientState")]),e._v(" is instantiated in the "),o("code",[e._v("InitGenesis")]),e._v(" handler of the 02-client submodule in core IBC.\nIt calls "),o("code",[e._v("CreateLocalhostClient")]),e._v(", declaring a new "),o("code",[e._v("ClientState")]),e._v(" and initializing it with its own client prefixed store.")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyAoayBLZWVwZXIpIENyZWF0ZUxvY2FsaG9zdENsaWVudChjdHggc2RrLkNvbnRleHQpIGVycm9yIHsKICB2YXIgY2xpZW50U3RhdGUgbG9jYWxob3N0LkNsaWVudFN0YXRlCiAgcmV0dXJuIGNsaWVudFN0YXRlLkluaXRpYWxpemUoY3R4LCBrLmNkYywgay5DbGllbnRTdG9yZShjdHgsIGV4cG9ydGVkLkxvY2FsaG9zdENsaWVudElEKSwgbmlsKQp9Cg=="}}),e._v(" "),o("p",[e._v("It is possible to disable the localhost client by removing the "),o("code",[e._v("09-localhost")]),e._v(" entry from the "),o("code",[e._v("allowed_clients")]),e._v(" list through governance.")]),e._v(" "),o("h2",{attrs:{id:"client-updates"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#client-updates"}},[e._v("#")]),e._v(" Client updates")]),e._v(" "),o("p",[e._v("The latest height is updated periodically through the ABCI "),o("a",{attrs:{href:"https://docs.cosmos.network/v0.47/building-modules/beginblock-endblock",target:"_blank",rel:"noopener noreferrer"}},[o("code",[e._v("BeginBlock")]),o("OutboundLink")],1),e._v(" interface of the 02-client submodule in core IBC.")]),e._v(" "),o("p",[o("a",{attrs:{href:"https://github.com/cosmos/ibc-go/blob/09-localhost/modules/core/02-client/abci.go#L12",target:"_blank",rel:"noopener noreferrer"}},[e._v("See "),o("code",[e._v("BeginBlocker")]),e._v(" in abci.go."),o("OutboundLink")],1)]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyBCZWdpbkJsb2NrZXIoY3R4IHNkay5Db250ZXh0LCBrIGtlZXBlci5LZWVwZXIpIHsKICAvLyAuLi4KCiAgaWYgY2xpZW50U3RhdGUsIGZvdW5kIDo9IGsuR2V0Q2xpZW50U3RhdGUoY3R4LCBleHBvcnRlZC5Mb2NhbGhvc3QpOyBmb3VuZCB7CiAgICBpZiBrLkdldENsaWVudFN0YXR1cyhjdHgsIGNsaWVudFN0YXRlLCBleHBvcnRlZC5Mb2NhbGhvc3QpID09IGV4cG9ydGVkLkFjdGl2ZSB7CiAgICAgIGsuVXBkYXRlTG9jYWxob3N0Q2xpZW50KGN0eCwgY2xpZW50U3RhdGUpCiAgICB9CiAgfQp9Cg=="}}),e._v(" "),o("p",[e._v("The above calls into the the 09-localhost "),o("code",[e._v("UpdateState")]),e._v(" method of the "),o("code",[e._v("ClientState")]),e._v(" .\nIt retrieves the current block height from the application context and sets the "),o("code",[e._v("LatestHeight")]),e._v(" of the 09-localhost client.")]),e._v(" "),o("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyAoY3MgQ2xpZW50U3RhdGUpIFVwZGF0ZVN0YXRlKGN0eCBzZGsuQ29udGV4dCwgY2RjIGNvZGVjLkJpbmFyeUNvZGVjLCBjbGllbnRTdG9yZSBzZGsuS1ZTdG9yZSwgY2xpZW50TXNnIGV4cG9ydGVkLkNsaWVudE1lc3NhZ2UpIFtdZXhwb3J0ZWQuSGVpZ2h0IHsKICBoZWlnaHQgOj0gY2xpZW50dHlwZXMuR2V0U2VsZkhlaWdodChjdHgpCiAgY3MuTGF0ZXN0SGVpZ2h0ID0gaGVpZ2h0CgogIGNsaWVudFN0b3JlLlNldChob3N0LkNsaWVudFN0YXRlS2V5KCksIGNsaWVudHR5cGVzLk11c3RNYXJzaGFsQ2xpZW50U3RhdGUoY2RjLCAmYW1wO2NzKSkKCiAgcmV0dXJuIFtdZXhwb3J0ZWQuSGVpZ2h0e2hlaWdodH0KfQo="}}),e._v(" "),o("p",[e._v("Note that the 09-localhost "),o("code",[e._v("ClientState")]),e._v(" is not updated through the 02-client interface leveraged by conventional IBC light clients.")])],1)}),[],!1,null,null,null);t.default=a.exports}}]);