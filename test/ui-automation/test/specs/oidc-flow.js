/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { wallet } = require("../helpers");

const vcSubjectData = [
  { name: "Given Name", value: "JOHN" },
  { name: "Family Name", value: "SMITH" },
  { name: "Gender", value: "Male" },
  { name: "Date of Birth", value: "1958-07-17" },
  { name: "Country of Birth", value: "Bahamas" },
  { name: "Resident Since", value: "2015-01-01" },
];


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
    await new Promise((resolve) => setTimeout(resolve, 3000));
  });

  it(`User validates the saved credential in Wallet`, async function () {
    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    const vcName = await $("span*=Permanent Resident Card");
    await vcName.waitForExist();
    await vcName.click();

    await wallet.validateCredentialDetails(vcSubjectData);
  });

  it(`User presents credential through OIDC Share: already signed-in`, async function () {
    // demo verifier page
    await browser.navigateTo(browser.config.oidcDemoVerifierURL);

    let oidcShareDemoBtn = await $("#oidc-share");
    await oidcShareDemoBtn.waitForExist();
    await oidcShareDemoBtn.click();

    const vcName = await $("span*=Permanent Resident Card");
    await vcName.waitForExist();

    await wallet.validateCredentialDetails(vcSubjectData);

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

