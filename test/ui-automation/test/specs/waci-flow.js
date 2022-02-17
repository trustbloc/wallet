/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { wallet } = require("../helpers");

const v2 = "v2";

const vcSubjectData = [
  { name: "Given Name", value: "JOHN" },
  { name: "Family Name", value: "SMITH" },
  { name: "Gender", value: "Male" },
  { name: "Date of birth", value: "1958-07-17" },
  { name: "Country of Birth", value: "Bahamas" },
  { name: "Resident Since", value: "2015-01-01" },
];

describe("TrustBloc Wallet - WACI flow", function () {
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

  describe(v2, function () {
    waciFlow(v2);
  });
});

async function waciFlow(version) {
  const ctx = {
    email: `ui-aut-${new Date().getTime()}@test.com`,
  };

  it(`User Sign up (${ctx.email})`, async function () {
    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.signUp(ctx);

    await wallet.waitForCredentials();

    // TODO - https://github.com/trustbloc/wallet/issues/1140 Dashboard loads before router connection is setup
    await new Promise((resolve) => setTimeout(resolve, 8000));
  });

  it(`User offered to save credential through WACI-Issuance (Redirect) : user (${ctx.email}) signed-in`, async function () {
    // demo issuer page
    await browser.navigateTo(browser.config.demoIssuerURL);

    let waciIssuanceDemoBtn;

    if (version === v2) {
      waciIssuanceDemoBtn = await $("#waci-issuance-demo-v2");
    }

    await waciIssuanceDemoBtn.waitForExist();
    await waciIssuanceDemoBtn.click();

    // accept store credential
    const storeButton = await $("#storeVCBtn");
    await storeButton.waitForClickable();
    await storeButton.click();

    const okBtn = await $("#issue-credentials-ok-btn");
    await okBtn.waitForExist();
    await okBtn.click();

    // success message
    const getSuccessMsg = await $("b*=Successfully Sent Credential to holder");
    await getSuccessMsg.waitForExist();
  });

  it(`User presents credential through WACI-Share (Redirect) : already signed-in`, async function () {
    // demo verifier page
    await browser.navigateTo(browser.config.demoVerifierURL);

    let waciShareDemoBtn;

    if (version === v2) {
      waciShareDemoBtn = await $("#waci-share-demo-v2");
    }

    await waciShareDemoBtn.waitForExist();
    await waciShareDemoBtn.click();

    const vcName = await $("span*=Permanent Resident Card");
    await vcName.waitForExist();

    await wallet.validateCredentialDetails(vcSubjectData);

    const shareBtn = await $("#share-credentials");
    await shareBtn.waitForExist();
    await shareBtn.click();

    const okBtn = await $("#share-credentials-ok-btn");
    await okBtn.waitForExist();
    await okBtn.click();

    const getSuccessMsg = await $("b*=Successfully Received Presentation");
    await getSuccessMsg.waitForExist();
  });

  it(`User signs out`, async function () {
    await browser.navigateTo(browser.config.walletURL);
    await wallet.signOut(ctx);
  });
}
