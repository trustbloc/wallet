/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const { wallet } = require("../helpers");

describe("TrustBloc Wallet - SignUp and SignIn flow", () => {
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
  });

  it(`Create Orb DID`, async function () {
    this.timeout(90000);

    await wallet.createOrbDID();
  });

  it(`Import DID Key (JWK key format)`, async function () {
    this.timeout(90000);

    await wallet.importDID({ method: 'key' });
  });

  // TODO add case to import DID with Base58 key format

  it(`Update Digital Identity preferences`, async function () {
    this.timeout(90000);

    await wallet.updatePreferences();
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
    await wallet.signIn(ctx);
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
