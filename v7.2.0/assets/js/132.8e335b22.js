(window.webpackJsonp=window.webpackJsonp||[]).push([[132],{690:function(e,t,a){"use strict";a.r(t);var n=a(1),o=Object(n.a)({},(function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[a("h1",{attrs:{id:"governance-proposals"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#governance-proposals"}},[e._v("#")]),e._v(" Governance Proposals")]),e._v(" "),a("p",[e._v("In uncommon situations, a highly valued client may become frozen due to uncontrollable\ncircumstances. A highly valued client might have hundreds of channels being actively used.\nSome of those channels might have a significant amount of locked tokens used for ICS 20.")]),e._v(" "),a("p",[e._v("If the one third of the validator set of the chain the client represents decides to collude,\nthey can sign off on two valid but conflicting headers each signed by the other one third\nof the honest validator set. The light client can now be updated with two valid, but conflicting\nheaders at the same height. The light client cannot know which header is trustworthy and therefore\nevidence of such misbehaviour is likely to be submitted resulting in a frozen light client.")]),e._v(" "),a("p",[e._v('Frozen light clients cannot be updated under any circumstance except via a governance proposal.\nSince a quorum of validators can sign arbitrary state roots which may not be valid executions\nof the state machine, a governance proposal has been added to ease the complexity of unfreezing\nor updating clients which have become "stuck". Without this mechanism, validator sets would need\nto construct a state root to unfreeze the client. Unfreezing clients, re-enables all of the channels\nbuilt upon that client. This may result in recovery of otherwise lost funds.')]),e._v(" "),a("p",[e._v("Tendermint light clients may become expired if the trusting period has passed since their\nlast update. This may occur if relayers stop submitting headers to update the clients.")]),e._v(" "),a("p",[e._v("An unplanned upgrade by the counterparty chain may also result in expired clients. If the counterparty\nchain undergoes an unplanned upgrade, there may be no commitment to that upgrade signed by the validator\nset before the chain-id changes. In this situation, the validator set of the last valid update for the\nlight client is never expected to produce another valid header since the chain-id has changed, which will\nultimately lead the on-chain light client to become expired.")]),e._v(" "),a("p",[e._v('In the case that a highly valued light client is frozen, expired, or rendered non-updateable, a\ngovernance proposal may be submitted to update this client, known as the subject client. The\nproposal includes the client identifier for the subject and the client identifier for a substitute\nclient. Light client implementations may implement custom updating logic, but in most cases,\nthe subject will be updated to the latest consensus state of the substitute client, if the proposal passes.\nThe substitute client is used as a "stand in" while the subject is on trial. It is best practice to create\na substitute client '),a("em",[e._v("after")]),e._v(" the subject has become frozen to avoid the substitute from also becoming frozen.\nAn active substitute client allows headers to be submitted during the voting period to prevent accidental expiry\nonce the proposal passes.")]),e._v(" "),a("p",[a("em",[e._v("note")]),e._v(" two of these parameters: "),a("code",[e._v("AllowUpdateAfterExpiry")]),e._v(" and "),a("code",[e._v("AllowUpdateAfterMisbehavior")]),e._v(" have been deprecated, and will both be set to "),a("code",[e._v("false")]),e._v(" upon upgrades even if they were previously set to "),a("code",[e._v("true")]),e._v(". These parameters will no longer play a role in restricting a client upgrade. Please see ADR026 for more details.")]),e._v(" "),a("h1",{attrs:{id:"how-to-recover-an-expired-client-with-a-governance-proposal"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#how-to-recover-an-expired-client-with-a-governance-proposal"}},[e._v("#")]),e._v(" How to recover an expired client with a governance proposal")]),e._v(" "),a("p",[e._v("See also the relevant documentation: "),a("RouterLink",{attrs:{to:"/architecture/adr-026-ibc-client-recovery-mechanisms.html"}},[e._v("ADR-026, IBC client recovery mechanisms")])],1),e._v(" "),a("blockquote",[a("p",[a("strong",[e._v("Who is this information for?")]),e._v("\nAlthough technically anyone can submit the governance proposal to recover an expired client, often it will be "),a("strong",[e._v("relayer operators")]),e._v(" (at least coordinating the submission).")])]),e._v(" "),a("h3",{attrs:{id:"preconditions"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#preconditions"}},[e._v("#")]),e._v(" Preconditions")]),e._v(" "),a("ul",[a("li",[e._v("The chain is updated with ibc-go >= v1.1.0.")]),e._v(" "),a("li",[e._v("There exists an active client (with a known client identifier) for the same counterparty chain as the expired client.")]),e._v(" "),a("li",[e._v("The governance deposit.")])]),e._v(" "),a("h2",{attrs:{id:"steps"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#steps"}},[e._v("#")]),e._v(" Steps")]),e._v(" "),a("h3",{attrs:{id:"step-1"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#step-1"}},[e._v("#")]),e._v(" Step 1")]),e._v(" "),a("p",[e._v("Check if the client is attached to the expected "),a("code",[e._v("chain-id")]),e._v(". For example, for an expired Tendermint client representing the Akash chain the client state looks like this on querying the client state:")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"ewogIGNsaWVudF9pZDogMDctdGVuZGVybWludC0xNDYKICBjbGllbnRfc3RhdGU6CiAgJ0B0eXBlJzogL2liYy5saWdodGNsaWVudHMudGVuZGVybWludC52MS5DbGllbnRTdGF0ZQogIGFsbG93X3VwZGF0ZV9hZnRlcl9leHBpcnk6IHRydWUKICBhbGxvd191cGRhdGVfYWZ0ZXJfbWlzYmVoYXZpb3VyOiB0cnVlCiAgY2hhaW5faWQ6IGFrYXNobmV0LTIKfQo="}}),e._v(" "),a("p",[e._v("The client is attached to the expected Akash "),a("code",[e._v("chain-id")]),e._v(". Note that although the parameters ("),a("code",[e._v("allow_update_after_expiry")]),e._v(" and "),a("code",[e._v("allow_update_after_misbehaviour")]),e._v(") exist to signal intent, these parameters have been deprecated and will not enforce any checks on the revival of client. See ADR-026 for more context on this deprecation.")]),e._v(" "),a("h3",{attrs:{id:"step-2"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#step-2"}},[e._v("#")]),e._v(" Step 2")]),e._v(" "),a("p",[e._v("If the chain has been updated to ibc-go >= v1.1.0, anyone can submit the governance proposal to recover the client by executing this via CLI.")]),e._v(" "),a("blockquote",[a("p",[e._v("Note that the Cosmos SDK has updated how governance proposals are submitted in SDK v0.46, now requiring to pass a .json proposal file")])]),e._v(" "),a("ul",[a("li",[a("p",[e._v("From SDK v0.46.x onwards")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"Jmx0O2JpbmFyeSZndDsgdHggZ292IHN1Ym1pdC1wcm9wb3NhbCBbcGF0aC10by1wcm9wb3NhbC1qc29uXQo="}}),e._v(" "),a("p",[e._v("where "),a("code",[e._v("proposal.json")]),e._v(" contains:")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"json",base64:"ewogICZxdW90O21lc3NhZ2VzJnF1b3Q7OiBbCiAgICB7CiAgICAgICZxdW90O0B0eXBlJnF1b3Q7OiAmcXVvdDsvaWJjLmNvcmUuY2xpZW50LnYxLkNsaWVudFVwZGF0ZVByb3Bvc2FsJnF1b3Q7LAogICAgICAmcXVvdDt0aXRsZSZxdW90OzogJnF1b3Q7dGl0bGVfc3RyaW5nJnF1b3Q7LAogICAgICAmcXVvdDtkZXNjcmlwdGlvbiZxdW90OzogJnF1b3Q7ZGVzY3JpcHRpb25fc3RyaW5nJnF1b3Q7LAogICAgICAmcXVvdDtzdWJqZWN0X2NsaWVudF9pZCZxdW90OzogJnF1b3Q7ZXhwaXJlZF9jbGllbnRfaWRfc3RyaW5nJnF1b3Q7LAogICAgICAmcXVvdDtzdWJzdGl0dXRlX2NsaWVudF9pZCZxdW90OzogJnF1b3Q7YWN0aXZlX2NsaWVudF9pZF9zdHJpbmcmcXVvdDsKICAgIH0KICBdLAogICZxdW90O21ldGFkYXRhJnF1b3Q7OiAmcXVvdDsmbHQ7bWV0YWRhdGEmZ3Q7JnF1b3Q7LAogICZxdW90O2RlcG9zaXQmcXVvdDs6ICZxdW90OzEwc3Rha2UmcXVvdDsKfQo="}}),e._v(" "),a("p",[e._v("Alternatively there's a legacy command (that is no longer recommended though):")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"Jmx0O2JpbmFyeSZndDsgdHggZ292IHN1Ym1pdC1sZWdhY3ktcHJvcG9zYWwgdXBkYXRlLWNsaWVudCAmbHQ7ZXhwaXJlZC1jbGllbnQtaWQmZ3Q7ICZsdDthY3RpdmUtY2xpZW50LWlkJmd0Owo="}})],1),e._v(" "),a("li",[a("p",[e._v("Until SDK v0.45.x")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"Jmx0O2JpbmFyeSZndDsgdHggZ292IHN1Ym1pdC1wcm9wb3NhbCB1cGRhdGUtY2xpZW50ICZsdDtleHBpcmVkLWNsaWVudC1pZCZndDsgJmx0O2FjdGl2ZS1jbGllbnQtaWQmZ3Q7Cg=="}})],1)]),e._v(" "),a("p",[e._v("The "),a("code",[e._v("<expired-client-id>")]),e._v(" identifier is the proposed client to be updated. This client must be either frozen or expired.")]),e._v(" "),a("p",[e._v("The "),a("code",[e._v("<active-client-id>")]),e._v(" represents a substitute client. It carries all the state for the client which may be updated. It must have identical client and chain parameters to the client which may be updated (except for latest height, frozen height, and chain ID). It should be continually updated during the voting period.")]),e._v(" "),a("p",[e._v("After this, all that remains is deciding who funds the governance deposit and ensuring the governance proposal passes. If it does, the client on trial will be updated to the latest state of the substitute.")]),e._v(" "),a("h2",{attrs:{id:"important-considerations"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#important-considerations"}},[e._v("#")]),e._v(" Important considerations")]),e._v(" "),a("p",[e._v("Please note that from v1.0.0 of ibc-go it will not be allowed for transactions to go to expired clients anymore, so please update to at least this version to prevent similar issues in the future.")]),e._v(" "),a("p",[e._v("Please also note that if the client on the other end of the transaction is also expired, that client will also need to update. This process updates only one client.")])],1)}),[],!1,null,null,null);t.default=o.exports}}]);