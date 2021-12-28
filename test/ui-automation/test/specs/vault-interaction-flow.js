/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import { wallet } from "../helpers";

describe("TrustBloc Wallet - Store/Share credential flow (CHAPI)", () => {
  const ctx = {
    email: `ui-aut-${new Date().getTime()}@test.com`,
  };

  // runs once before the first test in this block
  before(async () => {
    await browser.reloadSession();
    await browser.maximizeWindow();
  });

  beforeEach(function () {});

  it(`User Sign up (${ctx.email})`, async function () {
    this.timeout(90000);

    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.init(ctx);

    // 3. Go to Vaults screen
    const vaultsLink = await $("#navbar-link-vaults");
    await vaultsLink.click();
  });

  it(`Create Orb DID`, async function () {
    this.timeout(90000);

    await wallet.createOrbDID();
  });

  it(`Import DID Key (JWK key format)`, async function () {
    this.timeout(90000);

    await wallet.importDID({ method: "key" });
  });

  // TODO add case to import DID with Base58 key format

  it(`Update Digital Identity preferences`, async function () {
    this.timeout(90000);

    await wallet.updatePreferences();
  });

  it(`User performs DID Auth with mock issuer`, async function () {
    this.timeout(90000);

    // mock issuer (wallet page with sample requests)
    await browser.navigateTo(browser.config.webWalletURL);

    const didAuthBtn = await $("#didauth");
    await didAuthBtn.waitForExist();
    await didAuthBtn.click();

    const storeButton = await $("#get-btn");
    await storeButton.waitForClickable();
    await storeButton.click();

    await chapi.chooseWallet({
      name: browser.config.walletName,
    });

    await wallet.authenticate(ctx);
    await browser.switchToFrame(null);

    const didAuthSuccessMsg = await $(
      "div*=Successfully got response from wallet."
    );
    await didAuthSuccessMsg.waitForExist();
  });

  it(`User stores credential from mock issuer`, async function () {
    this.timeout(300000);

    for (const [key, value] of credential.entries()) {
      // mock issuer (wallet page with sample requests)
      await browser.navigateTo(browser.config.webWalletURL);

      console.log("save vc : start ", key);

      const vcSampleBtn = await $("#store-vc-sample-1");
      await vcSampleBtn.waitForExist();
      await vcSampleBtn.click();

      let vprs = {
        "@context": ["https://www.w3.org/2018/credentials/v1"],
        type: "VerifiablePresentation",
        verifiableCredential: [value.vc],
      };

      const sampleText = await $("#vcDataTextArea");
      await sampleText.clearValue();
      await sampleText.addValue(vprs);

      const storeButton = await $("#store-btn");
      await storeButton.waitForClickable();
      await storeButton.click();

      await chapi.chooseWallet({
        name: browser.config.walletName,
      });

      await wallet.validateCredentialDetails(value.vcSubjectData);

      await wallet.storeCredentials(ctx);
      await browser.switchToFrame(null);

      const storeSuccessMsg = await $(
        "div*=Successfully stored verifiable presentation to wallet."
      );
      await storeSuccessMsg.waitForExist();

      console.log("save vc : end ", key);
    }
  });

  it(`User validates the saved credential from mock issuer`, async function () {
    this.timeout(90000);

    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    for (const [key, value] of credential.entries()) {
      console.log("validate vc in wallet : start ", key);

      const vcName = await $("span*=" + value.name);
      await vcName.waitForExist();
      await vcName.click();

      await wallet.validateCredentialDetails(value.vcSubjectData);

      const credTab = await $("span*=Credentials");
      await credTab.waitForExist();
      await credTab.click();

      console.log("validate vc in wallet : end ", key);
    }
  });

  it(`User Sign Out (${ctx.email})`, async function () {
    this.timeout(90000);

    await wallet.logout(ctx);
  });

  it(`User Sign in (${ctx.email})`, async function () {
    this.timeout(90000);

    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Sign In to the registered Wallet (register/sign-up/etc.)
    await wallet.signIn(ctx.email);
  });

  it(`User validates the saved credential from mock issuer (after sign-in)`, async function () {
    this.timeout(90000);

    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    for (const [key, value] of credential.entries()) {
      console.log("validate vc in wallet : start ", key);

      const vcName = await $("span*=" + value.name);
      await vcName.waitForExist();
      await vcName.click();

      await wallet.validateCredentialDetails(value.vcSubjectData);

      const credTab = await $("span*=Credentials");
      await credTab.waitForExist();
      await credTab.click();

      console.log("validate vc in wallet : end ", key);
    }
  });

  it(`User shares the saved credential with mock verifier`, async function () {
    this.timeout(300000);

    for (const [key, value] of credential.entries()) {
      // mock verifier (wallet page with sample requests)
      await browser.navigateTo(browser.config.webWalletURL);

      console.log("share vc : start ", key);

      const vcSampleBtn = await $("#req-vp");
      await vcSampleBtn.waitForExist();
      await vcSampleBtn.click();

      let vprs = {
        web: {
          VerifiablePresentation: {
            query: [value.vpRequest],
            challenge: "4ada78c9-d58f-4e93-a039-cfa6c8999b97",
            domain: "example.com",
          },
        },
      };

      const sampleText = await $("#vcDataTextArea");
      await sampleText.clearValue();
      await sampleText.addValue(vprs);

      const getButton = await $("#get-btn");
      await getButton.waitForClickable();
      await getButton.click();

      await chapi.chooseWallet({
        name: browser.config.walletName,
      });

      // TODO https://github.com/trustbloc/wallet/issues/1124 VC Name mismatch between credentials screen and CHAPI share
      // const vcName = await $("span*=" + value.name);
      // await vcName.waitForExist();
      // await vcName.click();

      // await wallet.validateCredentialDetails(value.vcSubjectData);

      await wallet.presentCredentials(ctx);
      await browser.switchToFrame(null);

      const getSuccessMsg = await $(
        "div*=Successfully got response from wallet."
      );
      await getSuccessMsg.waitForExist();

      console.log("share vc : end ", key);
    }
  });
  it(`User deletes the saved credential`, async function () {
    this.timeout(90000);

    await browser.navigateTo(`${browser.config.walletURL}/credentials`);

    for (const [key, value] of credential.entries()) {
      console.log("validate vc in wallet : start ", key);

      const vcName = await $("span*=" + value.name);
      await vcName.waitForExist();
      await vcName.click();

      await wallet.validateCredentialDetails(value.vcSubjectData);
      console.log("validate vc in wallet : end ", key);
      await wallet.deleteCredential();
    }

    const isEmpty = credential.entries().length === 0;
    expect(isEmpty).toHaveValue("0");
  });

  it(`User Sign Out (${ctx.email})`, async function () {
    this.timeout(90000);

    await wallet.logout(ctx);
  });

  it(`User changes locale (${ctx.email})`, async function () {
    this.timeout(90000);

    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Change locale
    await wallet.changeLocale();
  });
});
