/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { wallet } = require("../helpers");

const prcvcSubjectData = [
  { name: "Given Name", value: "JOHN" },
  { name: "Family Name", value: "SMITH" },
  { name: "Gender", value: "Male" },
  { name: "Date of Birth", value: "1958-07-17" },
  { name: "Country of Birth", value: "Bahamas" },
  { name: "Resident Since", value: "2015-01-01" },
];

const udcvcSubjectData = [
  { name: "Degree", value: "Academic" },
  { name: "Type", value: "BachelorDegree" },
];

const issuerManifests = [
  {
    id: "GE-PRC-2022-CLASS-A",
    version: "0.1.0",
    issuer: {
      id: "did:example:123?linked-domains=3",
      name: "Government of Example Immigration",
      styles: {},
    },
    output_descriptors: [
      {
        id: "prc_output",
        schema: "https://w3id.org/citizenship/v1",
        display: {
          title: {
            path: ["$.name"],
            fallback: "Permanent Resident Card",
            schema: {
              type: "string",
            },
          },
          subtitle: {
            path: ["$.description"],
            fallback: "Government of Example Permanent Resident Card.",
            schema: {
              type: "string",
            },
          },
          description: {
            text: "Sample Permanent Resident Card issued by Government of Example Citizenship & Immigration Services",
          },
          properties: [
            {
              path: ["$.credentialSubject.image"],
              schema: {
                type: "string",
                contentMediaType: "image/png",
              },
              label: "Card Holder",
            },
            {
              path: ["$.credentialSubject.givenName"],
              schema: {
                type: "string",
              },
              label: "Given Name",
            },
            {
              path: ["$.credentialSubject.familyName"],
              schema: {
                type: "string",
              },
              label: "Family Name",
            },
            {
              path: ["$.credentialSubject.gender"],
              schema: {
                type: "string",
              },
              fallback: "Not disclosed",
              label: "Gender",
            },
            {
              path: ["$.credentialSubject.birthDate"],
              schema: {
                type: "string",
                format: "date",
              },
              label: "Date of Birth",
            },
            {
              path: ["$.credentialSubject.birthCountry"],
              schema: {
                type: "string",
              },
              label: "Country of Birth",
            },
            {
              path: ["$.credentialSubject.residentSince"],
              schema: {
                type: "string",
                format: "date",
              },
              label: "Resident Since",
            },
          ],
        },
        styles: {
          thumbnail: {
            uri: "https://file-server.trustbloc.local:12096/images/credential--uscis-icon.svg",
            alt: "Citizenship & Immigration Services",
          },
          hero: {
            uri: "https://example.com/trust.png",
            alt: "Service we trust",
          },
          background: {
            color: "#2b5283",
          },
          text: {
            color: "#fff",
          },
        },
      },
    ],
  },
  {
    id: "GE-UDC-2022",
    version: "0.1.0",
    issuer: {
      id: "did:example:123?linked-domains=3",
      name: "Example University",
      styles: {},
    },
    output_descriptors: [
      {
        id: "udc_output",
        schema: "https://www.w3.org/2018/credentials/examples/v1",
        display: {
          title: {
            path: ["$.name"],
            fallback: "University Degree Credential",
            schema: {
              type: "string",
            },
          },
          subtitle: {
            path: ["$.description"],
            fallback: "University of Example Degree.",
            schema: {
              type: "string",
            },
          },
          description: {
            text: "Sample University Degree issued by University of Example studies",
          },
          properties: [
            {
              path: ["$.credentialSubject.degree.name"],
              schema: {
                type: "string",
              },
              fallback: "Academic",
              label: "Degree",
            },
            {
              path: ["$.credentialSubject.degree.type"],
              schema: {
                type: "string",
              },
              fallback: "Not Specified",
              label: "Type",
            },
          ],
        },
        styles: {
          thumbnail: {
            uri: "credential--school-icon.svg",
            alt: "University of Example Studies",
          },
          hero: {
            uri: "https://example.com/happy-students.png",
            alt: "Happy Students",
          },
          background: {
            color: "#fff",
          },
          text: {
            color: "#190c21",
          },
        },
      },
    ],
  },
];

const credentials = {
  "https://w3id.org/citizenship/v1": {
    "@context": [
      "https://www.w3.org/2018/credentials/v1",
      "https://w3id.org/citizenship/v1",
    ],
    credentialSubject: {
      birthCountry: "Bahamas",
      birthDate: "1958-07-17",
      commuterClassification: "C1",
      familyName: "SMITH",
      gender: "Male",
      givenName: "JOHN",
      id: "did:example:b34ca6cd37bbf23",
      lprCategory: "C09",
      lprNumber: "999-999-999",
      residentSince: "2015-01-01",
      type: ["PermanentResident", "Person"],
    },
    description: "Government of Example Permanent Resident Card.",
    expirationDate: "2029-12-03T12:19:52Z",
    id: "https://issuer.oidp.uscis.gov/credentials/836274651",
    identifier: "83627465",
    issuanceDate: "2019-12-03T12:19:52Z",
    issuer: "did:example:b34ca6cd37bbf23",
    name: "Permanent Resident Card",
    type: ["VerifiableCredential", "PermanentResidentCard"],
  },
  "https://www.w3.org/2018/credentials/examples/v1": {
    "@context": [
      "https://www.w3.org/2018/credentials/v1",
      "https://www.w3.org/2018/credentials/examples/v1",
    ],
    credentialSchema: [],
    credentialSubject: {
      degree: {
        type: "BachelorDegree",
        university: "MIT",
      },
      id: "did:example:ebfeb1f712ebc6f1c276e12ec21",
      name: "Jayden Doe",
      spouse: "did:example:c276e12ec21ebfeb1f712ebc6f1",
    },

    name: "University Degree",
    description: "University Degree of Mr.John Smith",
    expirationDate: "2020-01-01T19:23:24Z",
    id: "http://example.edu/credentials/11873",
    issuanceDate: "2010-01-01T19:23:24Z",
    issuer: {
      id: "did:example:76e12ec712ebc6f1c221ebfeb1f",
      name: "Example University",
    },
    referenceNumber: 83294847,
    type: ["VerifiableCredential", "UniversityDegreeCredential"],
  },
};

describe("TrustBloc Wallet - OIDC flow", async function () {
  const ctx = {
    email: `ui-aut-oidc-${new Date().getTime()}@test.com`,
  };

  // runs once before the first test in this block
  before(async () => {
    await browser.reloadSession();
    await browser.maximizeWindow();
  });

  beforeEach(function () {});

  afterEach(async function () {
    if (this.currentTest.state === "failed") {
      const logs = await browser.getLogs("browser");
      console.log(JSON.stringify(logs, null, 4));
    }
  });

  it(`User Sign up (${ctx.email})`, async function () {
    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.signUp(ctx);

    await wallet.waitForCredentials();
  });

  it(`User offered to save credential through OIDC Issuance: already signed-in`, async function () {
    // demo issuer page
    await browser.navigateTo(browser.config.oidcDemoIssuerURL);

    let oidcIssuanceDemoBtn = await $("#oidc-issuance");
    await oidcIssuanceDemoBtn.waitForExist();
    await oidcIssuanceDemoBtn.click();

    const issuerLoginBtn = await $("#issuer-login");
    await issuerLoginBtn.waitForExist();
    await issuerLoginBtn.click();

    const prcSpan = await $("span*=Permanent Resident Card");
    await prcSpan.waitForExist();

    // accept store credential
    const storeButton = await $("#storeVCBtn");
    await storeButton.waitForClickable();
    await storeButton.click();

    const okBtn = await $("#issue-credentials-ok-btn");
    await okBtn.waitForExist();
    await okBtn.click();

    // sleep for 3 secs
    await new Promise((resolve) => setTimeout(resolve, 5000));
  });

  it(`User validates the saved credential in Wallet`, async function () {
    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    const vcName = await $("span*=Permanent Resident Card");
    await vcName.waitForExist();
    await vcName.click();

    await wallet.validateCredentialDetails(prcvcSubjectData);
  });

  it(`User presents credential through OIDC Share: already signed-in`, async function () {
    // demo verifier page
    await browser.navigateTo(browser.config.oidcDemoVerifierURL);

    let oidcShareDemoBtn = await $("#oidc-share");
    await oidcShareDemoBtn.waitForExist();
    await oidcShareDemoBtn.click();

    const vcName = await $("span*=Permanent Resident Card");
    await vcName.waitForExist();

    await wallet.validateCredentialDetails(prcvcSubjectData);

    const shareBtn = await $("#share-credentials");
    await shareBtn.waitForExist();
    await shareBtn.click();

    const msg = await $("b*=Successfully Received Presentation");
    await msg.waitForExist();
  });

  it(`User signs out`, async function () {
    await browser.navigateTo(browser.config.walletURL);
    await wallet.signOut(ctx);
  });
});

describe("TrustBloc Wallet - OIDC save multiple credential flow", async function () {
  const ctx = {
    email: `ui-aut-oidc-${new Date().getTime()}@test.com`,
  };

  // runs once before the first test in this block
  before(async () => {
    await browser.reloadSession();
    await browser.maximizeWindow();
  });

  beforeEach(function () {});

  afterEach(async function () {
    if (this.currentTest.state === "failed") {
      const logs = await browser.getLogs("browser");
      console.log(JSON.stringify(logs, null, 4));
    }
  });

  it(`User Sign up (${ctx.email})`, async function () {
    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.signUp(ctx);

    await wallet.waitForCredentials();
  });

  it(`User offered to save multiple credential through OIDC Issuance: already signed-in`, async function () {
    // demo issuer page
    await browser.navigateTo(browser.config.oidcDemoIssuerURL);

    let oidcIssuanceDemoBtn = await $("#oidc-issuance");
    await oidcIssuanceDemoBtn.waitForExist();

    let credentialTypes = await $("#credentialTypes");
    let manifestIDs = await $("#manifestIDs");
    let credentialManifests = await $("#credManifest");
    let credentialsToIssue = await $("#credsToIssue");

    await Promise.all([
      credentialTypes.waitForExist(),
      manifestIDs.waitForExist(),
      credentialManifests.waitForExist(),
      credentialsToIssue.waitForExist(),
    ]);

    await credentialTypes.setValue(
      "https://w3id.org/citizenship/v1,https://www.w3.org/2018/credentials/examples/v1"
    );
    await manifestIDs.setValue("GE-PRC-2022-CLASS-A,GE-UDC-2022");
    await credentialManifests.setValue(JSON.stringify(issuerManifests));
    await credentialsToIssue.setValue(JSON.stringify(credentials));

    await oidcIssuanceDemoBtn.click();

    const issuerLoginBtn = await $("#issuer-login");
    await issuerLoginBtn.waitForExist();
    await issuerLoginBtn.click();

    const prcSpan = await $("span*=Permanent Resident Card");
    await prcSpan.waitForExist();

    // accept store credential
    const storeButton = await $("#storeVCBtn");
    await storeButton.waitForClickable();
    await storeButton.click();

    const okBtn = await $("#issue-credentials-ok-btn");
    await okBtn.waitForExist();
    await okBtn.click();

    // sleep for 3 secs
    await new Promise((resolve) => setTimeout(resolve, 5000));
  });

  it(`User validates the saved credentials in Wallet`, async function () {
    await browser.navigateTo(`${browser.config.walletURL}/credentials`);
    const prcVC = await $("span*=Permanent Resident Card");
    await prcVC.waitForExist();
    await prcVC.click();
    await wallet.validateCredentialDetails(prcvcSubjectData);

    await browser.navigateTo(`${browser.config.walletURL}/credentials`);
    const udcVC = await $("span*=University Degree");
    await udcVC.waitForExist();
    await udcVC.click();
    await wallet.validateCredentialDetails(udcvcSubjectData);
  });

  it(`User signs out`, async function () {
    await browser.navigateTo(browser.config.walletURL);
    await wallet.signOut(ctx);
  });
});
