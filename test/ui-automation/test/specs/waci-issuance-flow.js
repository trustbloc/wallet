/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { wallet } = require("../helpers");

describe("TrustBloc Wallet - WACI Issuance flow", () => {
  const ctx = {
    email: `ui-aut-${new Date().getTime()}@test.com`,
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

    // TODO - https://github.com/trustbloc/wallet/issues/1140 Dashboard loads before router connection is setup
    await new Promise((resolve) => setTimeout(resolve, 8000));
  });

  it(`User offered to save credential through WACI-Issuance (Redirect) : user (${ctx.email}) signed-in`, async function () {
    // demo issuer page
    await browser.navigateTo(browser.config.demoIssuerURL);

    const waciIssuanceDemoBtn = await $("#waci-issuance-demo");
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

  it(`User signs out - (${ctx.email})`, async function () {
    await browser.navigateTo(browser.config.walletURL);
    await wallet.logout(ctx);
  });

  it(`User offered to save credential through WACI-Issuance (Redirect) : user (${ctx.email}) not signed-in`, async function () {
    // demo issuer page
    await browser.navigateTo(browser.config.demoIssuerURL);

    const waciIssuanceDemoBtn = await $("#waci-issuance-demo");
    await waciIssuanceDemoBtn.waitForExist();
    await waciIssuanceDemoBtn.click();

    await wallet.performSignIn(ctx.email, true);

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

  it(`User signs out`, async function () {
    await browser.navigateTo(browser.config.walletURL);
    await wallet.logout(ctx);
  });
});
