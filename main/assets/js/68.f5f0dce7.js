(window.webpackJsonp=window.webpackJsonp||[]).push([[68],{632:function(t,e,a){"use strict";a.r(e);var o=a(1),n=Object(o.a)({},(function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("ContentSlotsDistributor",{attrs:{"slot-key":t.$parent.slotKey}},[a("h1",{attrs:{id:"development-use-cases"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#development-use-cases"}},[t._v("#")]),t._v(" Development use cases")]),t._v(" "),a("p",[t._v("The initial version of Interchain Accounts allowed for the controller submodule to be extended by providing it with an underlying application which would handle all packet callbacks.\nThat functionality is now being deprecated in favor of alternative approaches.\nThis document will outline potential use cases and redirect each use case to the appropriate documentation.")]),t._v(" "),a("h2",{attrs:{id:"custom-authentication"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#custom-authentication"}},[t._v("#")]),t._v(" Custom authentication")]),t._v(" "),a("p",[t._v("Interchain accounts may be associated with alternative types of authentication relative to the traditional public/private key signing.\nIf you wish to develop or use Interchain Accounts with a custom authentication module and do not need to execute custom logic on the packet callbacks, we recommend you use ibc-go v6 or greater and that your custom authentication module interacts with the controller submodule via the "),a("RouterLink",{attrs:{to:"/apps/interchain-accounts/messages.html"}},[a("code",[t._v("MsgServer")])]),t._v(".")],1),t._v(" "),a("p",[t._v("If you wish to consume and execute custom logic in the packet callbacks, then please read the section "),a("a",{attrs:{href:"#packet-callbacks"}},[t._v("Packet callbacks")]),t._v(" below.")]),t._v(" "),a("h2",{attrs:{id:"redirection-to-a-smart-contract"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#redirection-to-a-smart-contract"}},[t._v("#")]),t._v(" Redirection to a smart contract")]),t._v(" "),a("p",[t._v("It may be desirable to allow smart contracts to control an interchain account.\nTo faciliate such an action, the controller submodule may be provided an underlying application which redirects to smart contract callers.\nAn improved design has been suggested in "),a("a",{attrs:{href:"https://github.com/cosmos/ibc-go/pull/1976",target:"_blank",rel:"noopener noreferrer"}},[t._v("ADR 008"),a("OutboundLink")],1),t._v(" which performs this action via middleware.")]),t._v(" "),a("p",[t._v("Implementors of this use case are recommended to follow the ADR 008 approach.\nThe underlying application may continue to be used as a short term solution for ADR 008 and the "),a("RouterLink",{attrs:{to:"/apps/interchain-accounts/auth-modules.html#registerinterchainaccount"}},[t._v("legacy API")]),t._v(" should continue to be utilized in such situations.")],1),t._v(" "),a("h2",{attrs:{id:"packet-callbacks"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#packet-callbacks"}},[t._v("#")]),t._v(" Packet callbacks")]),t._v(" "),a("p",[t._v("If a developer requires access to packet callbacks for their use case, then they have the following options:")]),t._v(" "),a("ol",[a("li",[t._v("Write a smart contract which is connected via an ADR 008 or equivalent IBC application (recommended).")]),t._v(" "),a("li",[t._v("Use the controller's underlying application to implement packet callback logic.")])]),t._v(" "),a("p",[t._v("In the first case, the smart contract should use the "),a("RouterLink",{attrs:{to:"/apps/interchain-accounts/messages.html"}},[a("code",[t._v("MsgServer")])]),t._v(".")],1),t._v(" "),a("p",[t._v("In the second case, the underlying application should use the "),a("RouterLink",{attrs:{to:"/apps/interchain-accounts/legacy/keeper-api.html"}},[t._v("legacy API")]),t._v(".")],1)])}),[],!1,null,null,null);e.default=n.exports}}]);