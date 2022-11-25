/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

import { wallet } from '../helpers';

const prcvcSubjectData = [
  { name: 'Given Name', value: 'JOHN' },
  { name: 'Family Name', value: 'SMITH' },
  { name: 'Gender', value: 'Male' },
  { name: 'Date of Birth', value: '1958-07-17' },
  { name: 'Country of Birth', value: 'Bahamas' },
  { name: 'Resident Since', value: '2015-01-01' },
];

describe('TrustBloc Wallet - OpenID4VC flow', function () {
  let ctx;

  // runs once before the first test in this block
  before(async function () {
    await browser.reloadSession();
    await browser.maximizeWindow();
    ctx = {
      email: `ui-aut-oidc-${new Date().getTime()}@test.com`,
    };
  });

  afterEach(async function () {
    if (this.currentTest.state === 'failed') {
      const logs = await browser.getLogs('browser');
      console.log(JSON.stringify(logs, null, 4));
    }
  });

  it(`User Sign up`, async function () {
    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.signUp(ctx, browser.config.isCHAPIEnabled);

    await wallet.waitForCredentials();

    // In this test, use JWT credentials instead of default JSON-LD, so both are tested.
    await wallet.useJWTCredentials();
  });

  it(`User offered to save credential through OIDC Issuance: already signed-in`, async function () {
    // demo issuer page
    await browser.navigateTo(browser.config.oidcDemoIssuerURL);

    const oidcIssuanceDemoBtn = await $('#oidc-issuance');
    await oidcIssuanceDemoBtn.waitForClickable();

    const walletUrlInput = await $('#walletURL');
    await walletUrlInput.waitForExist();
    await walletUrlInput.setValue(`${browser.config.walletURL}/oidc/initiate`);

    await oidcIssuanceDemoBtn.click();

    const issuerLoginBtn = await $('#issuer-login');
    await issuerLoginBtn.waitForClickable();
    await issuerLoginBtn.click();

    const prcSpan = await $('span*=Permanent Resident Card');
    await expect(prcSpan).toExist();

    // accept store credential
    const storeButton = await $('#storeVCBtn');
    await storeButton.waitForClickable();
    await storeButton.click();

    const okBtn = await $('#issue-credentials-ok-btn');
    await okBtn.waitForClickable();
    await okBtn.click();

    await expect(browser).toHaveUrl(`${browser.config.walletURL}/credentials`);
  });

  it(`User validates the saved credential in Wallet`, async function () {
    const vcName = await $('span*=Permanent Resident Card');
    await vcName.waitForClickable();
    await vcName.click();

    await wallet.validateCredentialDetails(prcvcSubjectData);
  });

  it(`User presents credential through OpenID4VP: already signed-in`, async function () {
    // demo verifier page
    await browser.navigateTo(`${browser.config.walletURL}/initiateProtocol`);

    const initiateInput = await $('#initiateIssuanceRequest');
    await initiateInput.waitForExist();
    await initiateInput.setValue(browser.config.openid4vpInitiateRequestURL);

    const initiateOpenID4VPButton = await $('#initiateFlow');
    await initiateOpenID4VPButton.waitForClickable();
    await initiateOpenID4VPButton.click();

    const vcName = await $('span*=Permanent Resident Card');
    await expect(vcName).toExist();

    await wallet.validateCredentialDetails(prcvcSubjectData);

    const shareBtn = await $('#share-credentials');
    await shareBtn.waitForClickable();
    await shareBtn.click();

    const successMsg = await $('#share-credentials-ok-btn');
    await expect(successMsg).toExist();
  });

  it(`User signs out`, async function () {
    await browser.navigateTo(browser.config.walletURL);
    await wallet.signOut(ctx);
  });
});
