(window.webpackJsonp=window.webpackJsonp||[]).push([[109],{668:function(e,t,a){"use strict";a.r(t);var o=a(1),i=Object(o.a)({},(function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[a("h1",{attrs:{id:"integration"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#integration"}},[e._v("#")]),e._v(" Integration")]),e._v(" "),a("p",{attrs:{synopsis:""}},[e._v("Learn how to integrate IBC to your application and send data packets to other chains.")]),e._v(" "),a("p",[e._v("This document outlines the required steps to integrate and configure the "),a("a",{attrs:{href:"https://github.com/cosmos/ibc-go/tree/main/modules/core",target:"_blank",rel:"noopener noreferrer"}},[e._v("IBC\nmodule"),a("OutboundLink")],1),e._v(" to your Cosmos SDK application and\nsend fungible token transfers to other chains.")]),e._v(" "),a("h2",{attrs:{id:"integrating-the-ibc-module"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#integrating-the-ibc-module"}},[e._v("#")]),e._v(" Integrating the IBC module")]),e._v(" "),a("p",[e._v("Integrating the IBC module to your SDK-based application is straighforward. The general changes can be summarized in the following steps:")]),e._v(" "),a("ul",[a("li",[e._v("Add required modules to the "),a("code",[e._v("module.BasicManager")])]),e._v(" "),a("li",[e._v("Define additional "),a("code",[e._v("Keeper")]),e._v(" fields for the new modules on the "),a("code",[e._v("App")]),e._v(" type")]),e._v(" "),a("li",[e._v("Add the module's "),a("code",[e._v("StoreKey")]),e._v("s and initialize their "),a("code",[e._v("Keeper")]),e._v("s")]),e._v(" "),a("li",[e._v("Set up corresponding routers and routes for the "),a("code",[e._v("ibc")]),e._v(" module")]),e._v(" "),a("li",[e._v("Add the modules to the module "),a("code",[e._v("Manager")])]),e._v(" "),a("li",[e._v("Add modules to "),a("code",[e._v("Begin/EndBlockers")]),e._v(" and "),a("code",[e._v("InitGenesis")])]),e._v(" "),a("li",[e._v("Update the module "),a("code",[e._v("SimulationManager")]),e._v(" to enable simulations")])]),e._v(" "),a("h3",{attrs:{id:"module-basicmanager-and-moduleaccount-permissions"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#module-basicmanager-and-moduleaccount-permissions"}},[e._v("#")]),e._v(" Module "),a("code",[e._v("BasicManager")]),e._v(" and "),a("code",[e._v("ModuleAccount")]),e._v(" permissions")]),e._v(" "),a("p",[e._v("The first step is to add the following modules to the "),a("code",[e._v("BasicManager")]),e._v(": "),a("code",[e._v("x/capability")]),e._v(", "),a("code",[e._v("x/ibc")]),e._v(",\nand "),a("code",[e._v("x/ibc-transfer")]),e._v(". After that, we need to grant "),a("code",[e._v("Minter")]),e._v(" and "),a("code",[e._v("Burner")]),e._v(" permissions to\nthe "),a("code",[e._v("ibc-transfer")]),e._v(" "),a("code",[e._v("ModuleAccount")]),e._v(" to mint and burn relayed tokens.")]),e._v(" "),a("h3",{attrs:{id:"integrating-light-clients"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#integrating-light-clients"}},[e._v("#")]),e._v(" Integrating light clients")]),e._v(" "),a("blockquote",[a("p",[e._v("Note that from v7 onwards, all light clients have to be explicitly registered in a chain's app.go and follow the steps listed below.\nThis is in contrast to earlier versions of ibc-go when "),a("code",[e._v("07-tendermint")]),e._v(" and "),a("code",[e._v("06-solomachine")]),e._v(" were added out of the box.")])]),e._v(" "),a("p",[e._v("All light clients must be registered with "),a("code",[e._v("module.BasicManager")]),e._v(" in a chain's app.go file.")]),e._v(" "),a("p",[e._v("The following code example shows how to register the existing "),a("code",[e._v("ibctm.AppModuleBasic{}")]),e._v(" light client implementation.")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"diff",base64:"CmltcG9ydCAoCiAgLi4uCisgaWJjdG0gJnF1b3Q7Z2l0aHViLmNvbS9jb3Ntb3MvaWJjLWdvL3Y2L21vZHVsZXMvbGlnaHQtY2xpZW50cy8wNy10ZW5kZXJtaW50JnF1b3Q7CiAgLi4uCikKCi8vIGFwcC5nbwp2YXIgKAoKICBNb2R1bGVCYXNpY3MgPSBtb2R1bGUuTmV3QmFzaWNNYW5hZ2VyKAogICAgLy8gLi4uCiAgICBjYXBhYmlsaXR5LkFwcE1vZHVsZUJhc2lje30sCiAgICBpYmMuQXBwTW9kdWxlQmFzaWN7fSwKICAgIHRyYW5zZmVyLkFwcE1vZHVsZUJhc2lje30sIC8vIGkuZSBpYmMtdHJhbnNmZXIgbW9kdWxlCgogICAgLy8gcmVnaXN0ZXIgbGlnaHQgY2xpZW50cyBvbiBJQkMKKyAgIGliY3RtLkFwcE1vZHVsZUJhc2lje30sCiAgKQoKICAvLyBtb2R1bGUgYWNjb3VudCBwZXJtaXNzaW9ucwogIG1hY2NQZXJtcyA9IG1hcFtzdHJpbmddW11zdHJpbmd7CiAgICAvLyBvdGhlciBtb2R1bGUgYWNjb3VudHMgcGVybWlzc2lvbnMKICAgIC8vIC4uLgogICAgaWJjdHJhbnNmZXJ0eXBlcy5Nb2R1bGVOYW1lOiAgICB7YXV0aHR5cGVzLk1pbnRlciwgYXV0aHR5cGVzLkJ1cm5lcn0sCiAgfQopCg=="}}),e._v(" "),a("h3",{attrs:{id:"application-fields"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#application-fields"}},[e._v("#")]),e._v(" Application fields")]),e._v(" "),a("p",[e._v("Then, we need to register the "),a("code",[e._v("Keepers")]),e._v(" as follows:")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"Ly8gYXBwLmdvCnR5cGUgQXBwIHN0cnVjdCB7CiAgLy8gYmFzZWFwcCwga2V5cyBhbmQgc3Vic3BhY2VzIGRlZmluaXRpb25zCgogIC8vIG90aGVyIGtlZXBlcnMKICAvLyAuLi4KICBJQkNLZWVwZXIgICAgICAgICppYmNrZWVwZXIuS2VlcGVyIC8vIElCQyBLZWVwZXIgbXVzdCBiZSBhIHBvaW50ZXIgaW4gdGhlIGFwcCwgc28gd2UgY2FuIFNldFJvdXRlciBvbiBpdCBjb3JyZWN0bHkKICBUcmFuc2ZlcktlZXBlciAgIGliY3RyYW5zZmVya2VlcGVyLktlZXBlciAvLyBmb3IgY3Jvc3MtY2hhaW4gZnVuZ2libGUgdG9rZW4gdHJhbnNmZXJzCgogIC8vIG1ha2Ugc2NvcGVkIGtlZXBlcnMgcHVibGljIGZvciB0ZXN0IHB1cnBvc2VzCiAgU2NvcGVkSUJDS2VlcGVyICAgICAgY2FwYWJpbGl0eWtlZXBlci5TY29wZWRLZWVwZXIKICBTY29wZWRUcmFuc2ZlcktlZXBlciBjYXBhYmlsaXR5a2VlcGVyLlNjb3BlZEtlZXBlcgoKICAvLy8gLi4uCiAgLy8vIG1vZHVsZSBhbmQgc2ltdWxhdGlvbiBtYW5hZ2VyIGRlZmluaXRpb25zCn0K"}}),e._v(" "),a("h3",{attrs:{id:"configure-the-keepers"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#configure-the-keepers"}},[e._v("#")]),e._v(" Configure the "),a("code",[e._v("Keepers")])]),e._v(" "),a("p",[e._v("During initialization, besides initializing the IBC "),a("code",[e._v("Keepers")]),e._v(" (for the "),a("code",[e._v("x/ibc")]),e._v(", and\n"),a("code",[e._v("x/ibc-transfer")]),e._v(" modules), we need to grant specific capabilities through the capability module\n"),a("code",[e._v("ScopedKeepers")]),e._v(" so that we can authenticate the object-capability permissions for each of the IBC\nchannels.")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"ZnVuYyBOZXdBcHAoLi4uYXJncykgKkFwcCB7CiAgLy8gZGVmaW5lIGNvZGVjcyBhbmQgYmFzZWFwcAoKICAvLyBhZGQgY2FwYWJpbGl0eSBrZWVwZXIgYW5kIFNjb3BlVG9Nb2R1bGUgZm9yIGliYyBtb2R1bGUKICBhcHAuQ2FwYWJpbGl0eUtlZXBlciA9IGNhcGFiaWxpdHlrZWVwZXIuTmV3S2VlcGVyKGFwcENvZGVjLCBrZXlzW2NhcGFiaWxpdHl0eXBlcy5TdG9yZUtleV0sIG1lbUtleXNbY2FwYWJpbGl0eXR5cGVzLk1lbVN0b3JlS2V5XSkKCiAgLy8gZ3JhbnQgY2FwYWJpbGl0aWVzIGZvciB0aGUgaWJjIGFuZCBpYmMtdHJhbnNmZXIgbW9kdWxlcwogIHNjb3BlZElCQ0tlZXBlciA6PSBhcHAuQ2FwYWJpbGl0eUtlZXBlci5TY29wZVRvTW9kdWxlKGliY2V4cG9ydGVkLk1vZHVsZU5hbWUpCiAgc2NvcGVkVHJhbnNmZXJLZWVwZXIgOj0gYXBwLkNhcGFiaWxpdHlLZWVwZXIuU2NvcGVUb01vZHVsZShpYmN0cmFuc2ZlcnR5cGVzLk1vZHVsZU5hbWUpCgogIC8vIC4uLiBvdGhlciBtb2R1bGVzIGtlZXBlcnMKCiAgLy8gQ3JlYXRlIElCQyBLZWVwZXIKICBhcHAuSUJDS2VlcGVyID0gaWJja2VlcGVyLk5ld0tlZXBlcigKICAgIGFwcENvZGVjLCBrZXlzW2liY2V4cG9ydGVkLlN0b3JlS2V5XSwgYXBwLkdldFN1YnNwYWNlKGliY2V4cG9ydGVkLk1vZHVsZU5hbWUpLCBhcHAuU3Rha2luZ0tlZXBlciwgYXBwLlVwZ3JhZGVLZWVwZXIsIHNjb3BlZElCQ0tlZXBlciwKICApCgogIC8vIENyZWF0ZSBUcmFuc2ZlciBLZWVwZXJzCiAgYXBwLlRyYW5zZmVyS2VlcGVyID0gaWJjdHJhbnNmZXJrZWVwZXIuTmV3S2VlcGVyKAogICAgYXBwQ29kZWMsIGtleXNbaWJjdHJhbnNmZXJ0eXBlcy5TdG9yZUtleV0sIGFwcC5HZXRTdWJzcGFjZShpYmN0cmFuc2ZlcnR5cGVzLk1vZHVsZU5hbWUpLAogICAgYXBwLklCQ0tlZXBlci5DaGFubmVsS2VlcGVyLCBhcHAuSUJDS2VlcGVyLkNoYW5uZWxLZWVwZXIsICZhbXA7YXBwLklCQ0tlZXBlci5Qb3J0S2VlcGVyLAogICAgYXBwLkFjY291bnRLZWVwZXIsIGFwcC5CYW5rS2VlcGVyLCBzY29wZWRUcmFuc2ZlcktlZXBlciwKICApCiAgdHJhbnNmZXJNb2R1bGUgOj0gdHJhbnNmZXIuTmV3QXBwTW9kdWxlKGFwcC5UcmFuc2ZlcktlZXBlcikKCiAgLy8gLi4gY29udGludWVzCn0K"}}),e._v(" "),a("h3",{attrs:{id:"register-routers"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#register-routers"}},[e._v("#")]),e._v(" Register "),a("code",[e._v("Routers")])]),e._v(" "),a("p",[e._v("IBC needs to know which module is bound to which port so that it can route packets to the\nappropriate module and call the appropriate callbacks. The port to module name mapping is handled by\nIBC's port "),a("code",[e._v("Keeper")]),e._v(". However, the mapping from module name to the relevant callbacks is accomplished\nby the port\n"),a("a",{attrs:{href:"https://github.com/cosmos/ibc-go/blob/main/modules/core/05-port/types/router.go",target:"_blank",rel:"noopener noreferrer"}},[a("code",[e._v("Router")]),a("OutboundLink")],1),e._v(" on the\nIBC module.")]),e._v(" "),a("p",[e._v("Adding the module routes allows the IBC handler to call the appropriate callback when processing a\nchannel handshake or a packet.")]),e._v(" "),a("p",[e._v("Currently, a "),a("code",[e._v("Router")]),e._v(" is static so it must be initialized and set correctly on app initialization.\nOnce the "),a("code",[e._v("Router")]),e._v(" has been set, no new routes can be added.")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"Ly8gYXBwLmdvCmZ1bmMgTmV3QXBwKC4uLmFyZ3MpICpBcHAgewogIC8vIC4uIGNvbnRpbnVhdGlvbiBmcm9tIGFib3ZlCgogIC8vIENyZWF0ZSBzdGF0aWMgSUJDIHJvdXRlciwgYWRkIGliYy10cmFuZmVyIG1vZHVsZSByb3V0ZSwgdGhlbiBzZXQgYW5kIHNlYWwgaXQKICBpYmNSb3V0ZXIgOj0gcG9ydC5OZXdSb3V0ZXIoKQogIGliY1JvdXRlci5BZGRSb3V0ZShpYmN0cmFuc2ZlcnR5cGVzLk1vZHVsZU5hbWUsIHRyYW5zZmVyTW9kdWxlKQogIC8vIFNldHRpbmcgUm91dGVyIHdpbGwgZmluYWxpemUgYWxsIHJvdXRlcyBieSBzZWFsaW5nIHJvdXRlcgogIC8vIE5vIG1vcmUgcm91dGVzIGNhbiBiZSBhZGRlZAogIGFwcC5JQkNLZWVwZXIuU2V0Um91dGVyKGliY1JvdXRlcikKCiAgLy8gLi4gY29udGludWVzCg=="}}),e._v(" "),a("h3",{attrs:{id:"module-managers"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#module-managers"}},[e._v("#")]),e._v(" Module Managers")]),e._v(" "),a("p",[e._v("In order to use IBC, we need to add the new modules to the module "),a("code",[e._v("Manager")]),e._v(" and to the "),a("code",[e._v("SimulationManager")]),e._v(" in case your application supports "),a("a",{attrs:{href:"https://github.com/cosmos/cosmos-sdk/blob/main/docs/docs/building-modules/14-simulator.md",target:"_blank",rel:"noopener noreferrer"}},[e._v("simulations"),a("OutboundLink")],1),e._v(".")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"Ly8gYXBwLmdvCmZ1bmMgTmV3QXBwKC4uLmFyZ3MpICpBcHAgewogIC8vIC4uIGNvbnRpbnVhdGlvbiBmcm9tIGFib3ZlCgogIGFwcC5tbSA9IG1vZHVsZS5OZXdNYW5hZ2VyKAogICAgLy8gb3RoZXIgbW9kdWxlcwogICAgLy8gLi4uCiAgICBjYXBhYmlsaXR5Lk5ld0FwcE1vZHVsZShhcHBDb2RlYywgKmFwcC5DYXBhYmlsaXR5S2VlcGVyKSwKICAgIGliYy5OZXdBcHBNb2R1bGUoYXBwLklCQ0tlZXBlciksCiAgICB0cmFuc2Zlck1vZHVsZSwKICApCgogIC8vIC4uLgoKICBhcHAuc20gPSBtb2R1bGUuTmV3U2ltdWxhdGlvbk1hbmFnZXIoCiAgICAvLyBvdGhlciBtb2R1bGVzCiAgICAvLyAuLi4KICAgIGNhcGFiaWxpdHkuTmV3QXBwTW9kdWxlKGFwcENvZGVjLCAqYXBwLkNhcGFiaWxpdHlLZWVwZXIpLAogICAgaWJjLk5ld0FwcE1vZHVsZShhcHAuSUJDS2VlcGVyKSwKICAgIHRyYW5zZmVyTW9kdWxlLAogICkKCiAgLy8gLi4gY29udGludWVzCg=="}}),e._v(" "),a("h3",{attrs:{id:"application-abci-ordering"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#application-abci-ordering"}},[e._v("#")]),e._v(" Application ABCI Ordering")]),e._v(" "),a("p",[e._v("One addition from IBC is the concept of "),a("code",[e._v("HistoricalEntries")]),e._v(" which are stored on the staking module.\nEach entry contains the historical information for the "),a("code",[e._v("Header")]),e._v(" and "),a("code",[e._v("ValidatorSet")]),e._v(" of this chain which is stored\nat each height during the "),a("code",[e._v("BeginBlock")]),e._v(" call. The historical info is required to introspect the\npast historical info at any given height in order to verify the light client "),a("code",[e._v("ConsensusState")]),e._v(" during the\nconnection handhake.")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"go",base64:"Ly8gYXBwLmdvCmZ1bmMgTmV3QXBwKC4uLmFyZ3MpICpBcHAgewogIC8vIC4uIGNvbnRpbnVhdGlvbiBmcm9tIGFib3ZlCgogIC8vIGFkZCBzdGFraW5nIGFuZCBpYmMgbW9kdWxlcyB0byBCZWdpbkJsb2NrZXJzCiAgYXBwLm1tLlNldE9yZGVyQmVnaW5CbG9ja2VycygKICAgIC8vIG90aGVyIG1vZHVsZXMgLi4uCiAgICBzdGFraW5ndHlwZXMuTW9kdWxlTmFtZSwgaWJjZXhwb3J0ZWQuTW9kdWxlTmFtZSwKICApCgogIC8vIC4uLgoKICAvLyBOT1RFOiBDYXBhYmlsaXR5IG1vZHVsZSBtdXN0IG9jY3VyIGZpcnN0IHNvIHRoYXQgaXQgY2FuIGluaXRpYWxpemUgYW55IGNhcGFiaWxpdGllcwogIC8vIHNvIHRoYXQgb3RoZXIgbW9kdWxlcyB0aGF0IHdhbnQgdG8gY3JlYXRlIG9yIGNsYWltIGNhcGFiaWxpdGllcyBhZnRlcndhcmRzIGluIEluaXRDaGFpbgogIC8vIGNhbiBkbyBzbyBzYWZlbHkuCiAgYXBwLm1tLlNldE9yZGVySW5pdEdlbmVzaXMoCiAgICBjYXBhYmlsaXR5dHlwZXMuTW9kdWxlTmFtZSwKICAgIC8vIG90aGVyIG1vZHVsZXMgLi4uCiAgICBpYmNleHBvcnRlZC5Nb2R1bGVOYW1lLCBpYmN0cmFuc2ZlcnR5cGVzLk1vZHVsZU5hbWUsCiAgKQoKICAvLyAuLiBjb250aW51ZXMK"}}),e._v(" "),a("div",{staticClass:"custom-block warning"},[a("p",[a("strong",[e._v("IMPORTANT")]),e._v(": The capability module "),a("strong",[e._v("must")]),e._v(" be declared first in "),a("code",[e._v("SetOrderInitGenesis")])])]),e._v(" "),a("p",[e._v("That's it! You have now wired up the IBC module and are now able to send fungible tokens across\ndifferent chains. If you want to have a broader view of the changes take a look into the SDK's\n"),a("a",{attrs:{href:"https://github.com/cosmos/ibc-go/blob/main/testing/simapp/app.go",target:"_blank",rel:"noopener noreferrer"}},[a("code",[e._v("SimApp")]),a("OutboundLink")],1),e._v(".")]),e._v(" "),a("h2",{attrs:{hide:"",id:"next"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#next"}},[e._v("#")]),e._v(" Next")]),e._v(" "),a("p",{attrs:{hide:""}},[e._v("Learn about how to create "),a("RouterLink",{attrs:{to:"/ibc/apps/apps.html"}},[e._v("custom IBC modules")]),e._v(" for your application")],1)],1)}),[],!1,null,null,null);t.default=i.exports}}]);