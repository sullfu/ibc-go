(window.webpackJsonp=window.webpackJsonp||[]).push([[94],{655:function(e,t,n){"use strict";n.r(t);var a=n(1),i=Object(a.a)({},(function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[n("h1",{attrs:{id:"adr-27-add-support-for-wasm-based-light-client"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#adr-27-add-support-for-wasm-based-light-client"}},[e._v("#")]),e._v(" ADR 27: Add support for Wasm based light client")]),e._v(" "),n("h2",{attrs:{id:"changelog"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#changelog"}},[e._v("#")]),e._v(" Changelog")]),e._v(" "),n("ul",[n("li",[e._v("26/11/2020: Initial Draft")])]),e._v(" "),n("h2",{attrs:{id:"status"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#status"}},[e._v("#")]),e._v(" Status")]),e._v(" "),n("p",[n("em",[e._v("Draft, needs updates")])]),e._v(" "),n("h2",{attrs:{id:"abstract"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#abstract"}},[e._v("#")]),e._v(" Abstract")]),e._v(" "),n("p",[e._v("In the Cosmos SDK light clients are current hardcoded in Go. This makes upgrading existing IBC light clients or adding\nsupport for new light client a multi step process involving on-chain governance which is time-consuming.")]),e._v(" "),n("p",[e._v("To remedy this, we are proposing a WASM VM to host light client bytecode, which allows easier upgrading of\nexisting IBC light clients as well as adding support for new IBC light clients without requiring a code release and corresponding\nhard-fork event.")]),e._v(" "),n("h2",{attrs:{id:"context"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#context"}},[e._v("#")]),e._v(" Context")]),e._v(" "),n("p",[e._v("Currently in the SDK, light clients are defined as part of the codebase and are implemented as submodules under\n"),n("code",[e._v("ibc-go/core/modules/light-clients/")]),e._v(".")]),e._v(" "),n("p",[e._v("Adding support for new light client or update an existing light client in the event of security\nissue or consensus update is multi-step process which is both time consuming and error prone:")]),e._v(" "),n("ol",[n("li",[n("p",[e._v("To add support for new light client or update an existing light client in the\nevent of security issue or consensus update, we need to modify the codebase and integrate it in numerous places.")])]),e._v(" "),n("li",[n("p",[e._v("Governance voting: Adding new light client implementations require governance support and is expensive: This is\nnot ideal as chain governance is gatekeeper for new light client implementations getting added. If a small community\nwant support for light client X, they may not be able to convince governance to support it.")])]),e._v(" "),n("li",[n("p",[e._v("Validator upgrade: After governance voting succeeds, validators need to upgrade their nodes in order to enable new\nIBC light client implementation.")])])]),e._v(" "),n("p",[e._v("Another problem stemming from the above process is that if a chain wants to upgrade its own consensus, it will need to convince every chain\nor hub connected to it to upgrade its light client in order to stay connected. Due to time consuming process required\nto upgrade light client, a chain with lots of connections needs to be disconnected for quite some time after upgrading\nits consensus, which can be very expensive in terms of time and effort.")]),e._v(" "),n("p",[e._v("We are proposing simplifying this workflow by integrating a WASM light client module which makes adding support for\na new light client a simple transaction. The light client bytecode, written in Wasm-compilable Rust, runs inside a WASM\nVM. The Wasm light client submodule exposes a proxy light client interface that routes incoming messages to the\nappropriate handler function, inside the Wasm VM for execution.")]),e._v(" "),n("p",[e._v("With WASM light client module, anybody can add new IBC light client in the form of WASM bytecode (provided they are able to pay the requisite gas fee for the transaction)\nas well as instantiate clients using any created client type. This allows any chain to update its own light client in other chains\nwithout going through steps outlined above.")]),e._v(" "),n("h2",{attrs:{id:"decision"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#decision"}},[e._v("#")]),e._v(" Decision")]),e._v(" "),n("p",[e._v("We decided to use WASM light client module as a light client proxy which will interface with the actual light client\nuploaded as WASM bytecode. This will require changing client selection method to allow any client if the client type\nhas prefix of "),n("code",[e._v("wasm/")]),e._v(".")]),e._v(" "),n("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"Ly8gSXNBbGxvd2VkQ2xpZW50IGNoZWNrcyBpZiB0aGUgZ2l2ZW4gY2xpZW50IHR5cGUgaXMgcmVnaXN0ZXJlZCBvbiB0aGUgYWxsb3dsaXN0LgpmdW5jIChwIFBhcmFtcykgSXNBbGxvd2VkQ2xpZW50KGNsaWVudFR5cGUgc3RyaW5nKSBib29sIHsKCWlmIHAuQXJlV0FTTUNsaWVudHNBbGxvd2VkICZhbXA7JmFtcDsgaXNXQVNNQ2xpZW50KGNsaWVudFR5cGUpIHsKCQlyZXR1cm4gdHJ1ZQoJfQoJCglmb3IgXywgYWxsb3dlZENsaWVudCA6PSByYW5nZSBwLkFsbG93ZWRDbGllbnRzIHsKCQlpZiBhbGxvd2VkQ2xpZW50ID09IGNsaWVudFR5cGUgewoJCQlyZXR1cm4gdHJ1ZQoJCX0KCX0KCglyZXR1cm4gZmFsc2UKfQo="}}),e._v(" "),n("p",[e._v("To upload new light client, user need to create a transaction with Wasm byte code which will be\nprocessed by IBC Wasm module.")]),e._v(" "),n("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyAoayBLZWVwZXIpIFVwbG9hZExpZ2h0Q2xpZW50ICh3YXNtQ29kZTogW11ieXRlLCBkZXNjcmlwdGlvbjogU3RyaW5nKSB7CiAgICB3YXNtUmVnaXN0cnkgPSBnZXRXQVNNUmVnaXN0cnkoKQogICAgaWQgOj0gaGV4LkVuY29kZVRvU3RyaW5nKHNoYTI1Ni5TdW0yNTYod2FzbUNvZGUpKQogICAgYXNzZXJ0KCF3YXNtUmVnaXN0cnkuRXhpc3RzKGlkKSkKICAgIGFzc2VydCh3YXNtUmVnaXN0cnkuVmFsaWRhdGVBbmRTdG9yZUNvZGUoaWQsIGRlc2NyaXB0aW9uLCB3YXNtQ29kZSwgZmFsc2UpKQp9Cg=="}}),e._v(" "),n("p",[e._v("As name implies, Wasm registry is a registry which stores set of Wasm client code indexed by its hash and allows\nclient code to retrieve latest code uploaded.")]),e._v(" "),n("p",[n("code",[e._v("ValidateAndStoreCode")]),e._v(" checks if the wasm bytecode uploaded is valid and confirms to VM interface.")]),e._v(" "),n("h3",{attrs:{id:"how-light-client-proxy-works"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#how-light-client-proxy-works"}},[e._v("#")]),e._v(" How light client proxy works?")]),e._v(" "),n("p",[e._v("The light client proxy behind the scenes will call a cosmwasm smart contract instance with incoming arguments in json\nserialized format with appropriate environment information. Data returned by the smart contract is deserialized and\nreturned to the caller.")]),e._v(" "),n("p",[e._v("Consider an example of "),n("code",[e._v("CheckProposedHeaderAndUpdateState")]),e._v(" function of "),n("code",[e._v("ClientState")]),e._v(" interface. Incoming arguments are\npackaged inside a payload which is json serialized and passed to "),n("code",[e._v("callContract")]),e._v(" which calls "),n("code",[e._v("vm.Execute")]),e._v(" and returns the\narray of bytes returned by the smart contract. This data is deserialized and passed as return argument.")]),e._v(" "),n("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyAoYyAqQ2xpZW50U3RhdGUpIENoZWNrUHJvcG9zZWRIZWFkZXJBbmRVcGRhdGVTdGF0ZShjb250ZXh0IHNkay5Db250ZXh0LCBtYXJzaGFsZXIgY29kZWMuQmluYXJ5TWFyc2hhbGVyLCBzdG9yZSBzZGsuS1ZTdG9yZSwgaGVhZGVyIGV4cG9ydGVkLkNsaWVudE1lc3NhZ2UpIChleHBvcnRlZC5DbGllbnRTdGF0ZSwgZXhwb3J0ZWQuQ29uc2Vuc3VzU3RhdGUsIGVycm9yKSB7CgkvLyBnZXQgY29uc2Vuc3VzIHN0YXRlIGNvcnJlc3BvbmRpbmcgdG8gY2xpZW50IHN0YXRlIHRvIGNoZWNrIGlmIHRoZSBjbGllbnQgaXMgZXhwaXJlZAoJY29uc2Vuc3VzU3RhdGUsIGVyciA6PSBHZXRDb25zZW5zdXNTdGF0ZShzdG9yZSwgbWFyc2hhbGVyLCBjLkxhdGVzdEhlaWdodCkKCWlmIGVyciAhPSBuaWwgewoJCXJldHVybiBuaWwsIG5pbCwgc2RrZXJyb3JzLldyYXBmKAoJCQllcnIsICZxdW90O2NvdWxkIG5vdCBnZXQgY29uc2Vuc3VzIHN0YXRlIGZyb20gY2xpZW50c3RvcmUgYXQgaGVpZ2h0OiAlZCZxdW90OywgYy5MYXRlc3RIZWlnaHQsCgkJKQoJfQoJCglwYXlsb2FkIDo9IG1ha2UobWFwW3N0cmluZ11tYXBbc3RyaW5nXWludGVyZmFjZXt9KQoJcGF5bG9hZFtDaGVja1Byb3Bvc2VkSGVhZGVyQW5kVXBkYXRlU3RhdGVdID0gbWFrZShtYXBbc3RyaW5nXWludGVyZmFjZXt9KQoJaW5uZXIgOj0gcGF5bG9hZFtDaGVja1Byb3Bvc2VkSGVhZGVyQW5kVXBkYXRlU3RhdGVdCglpbm5lclsmcXVvdDttZSZxdW90O10gPSBjCglpbm5lclsmcXVvdDtoZWFkZXImcXVvdDtdID0gaGVhZGVyCglpbm5lclsmcXVvdDtjb25zZW5zdXNfc3RhdGUmcXVvdDtdID0gY29uc2Vuc3VzU3RhdGUKCgllbmNvZGVkRGF0YSwgZXJyIDo9IGpzb24uTWFyc2hhbChwYXlsb2FkKQoJaWYgZXJyICE9IG5pbCB7CgkJcmV0dXJuIG5pbCwgbmlsLCBzZGtlcnJvcnMuV3JhcGYoRXJyVW5hYmxlVG9NYXJzaGFsUGF5bG9hZCwgZm10LlNwcmludGYoJnF1b3Q7dW5kZXJseWluZyBlcnJvcjogJXMmcXVvdDssIGVyci5FcnJvcigpKSkKCX0KCW91dCwgZXJyIDo9IGNhbGxDb250cmFjdChjLkNvZGVJZCwgY29udGV4dCwgc3RvcmUsIGVuY29kZWREYXRhKQoJaWYgZXJyICE9IG5pbCB7CgkJcmV0dXJuIG5pbCwgbmlsLCBzZGtlcnJvcnMuV3JhcGYoRXJyVW5hYmxlVG9DYWxsLCBmbXQuU3ByaW50ZigmcXVvdDt1bmRlcmx5aW5nIGVycm9yOiAlcyZxdW90OywgZXJyLkVycm9yKCkpKQoJfQoJb3V0cHV0IDo9IGNsaWVudFN0YXRlQ2FsbFJlc3BvbnNle30KCWlmIGVyciA6PSBqc29uLlVubWFyc2hhbChvdXQuRGF0YSwgJmFtcDtvdXRwdXQpOyBlcnIgIT0gbmlsIHsKCQlyZXR1cm4gbmlsLCBuaWwsIHNka2Vycm9ycy5XcmFwZihFcnJVbmFibGVUb1VubWFyc2hhbFBheWxvYWQsIGZtdC5TcHJpbnRmKCZxdW90O3VuZGVybHlpbmcgZXJyb3I6ICVzJnF1b3Q7LCBlcnIuRXJyb3IoKSkpCgl9CglpZiAhb3V0cHV0LlJlc3VsdC5Jc1ZhbGlkIHsKCQlyZXR1cm4gbmlsLCBuaWwsIGZtdC5FcnJvcmYoJnF1b3Q7JXMgZXJyb3Igb2N1cnJlZCB3aGlsZSB1cGRhdGluZyBjbGllbnQgc3RhdGUmcXVvdDssIG91dHB1dC5SZXN1bHQuRXJyb3JNc2cpCgl9CglvdXRwdXQucmVzZXRJbW11dGFibGVzKGMpCglyZXR1cm4gb3V0cHV0Lk5ld0NsaWVudFN0YXRlLCBvdXRwdXQuTmV3Q29uc2Vuc3VzU3RhdGUsIG5pbAp9Cg=="}}),e._v(" "),n("h2",{attrs:{id:"consequences"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#consequences"}},[e._v("#")]),e._v(" Consequences")]),e._v(" "),n("h3",{attrs:{id:"positive"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#positive"}},[e._v("#")]),e._v(" Positive")]),e._v(" "),n("ul",[n("li",[e._v("Adding support for new light client or upgrading existing light client is way easier than before and only requires single transaction.")]),e._v(" "),n("li",[e._v("Improves maintainability of Cosmos SDK, since no change in codebase is required to support new client or upgrade it.")])]),e._v(" "),n("h3",{attrs:{id:"negative"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#negative"}},[e._v("#")]),e._v(" Negative")]),e._v(" "),n("ul",[n("li",[e._v("Light clients need to be written in subset of rust which could compile in Wasm.")]),e._v(" "),n("li",[e._v("Introspecting light client code is difficult as only compiled bytecode exists in the blockchain.")])])],1)}),[],!1,null,null,null);t.default=i.exports}}]);