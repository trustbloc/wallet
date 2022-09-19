/*
Copyright Digital Bazaar

This file was copied from https://github.com/w3c-ccg/chapi-interop-test-suite/blob/main/test/helpers/chapi.js.
The license details are available at https://github.com/w3c-ccg/chapi-interop-test-suite#license.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

exports.chooseWallet = async ({ name }) => {
  const chapiFrame = await $('iframe');
  await chapiFrame.waitForExist();
  await browser.switchToFrame(chapiFrame);

  const rememberChoiceBtn = await $('span*=Remember');
  await rememberChoiceBtn.waitForClickable();
  await rememberChoiceBtn.click();

  const demoWallet = await $(`span*=${name}`);
  await demoWallet.waitForClickable();
  await demoWallet.click();

  const innerWalletFrame = await $('iframe');
  await innerWalletFrame.waitForExist();
  await browser.switchToFrame(innerWalletFrame);

  let dialogs;
  await browser.waitUntil(async () => {
    dialogs = await $$('dialog');
    return dialogs.length === 2;
  });

  const innerWalletFrame2 = await dialogs[1].$('iframe');
  await innerWalletFrame2.waitForExist();
  await browser.switchToFrame(innerWalletFrame2);
};

exports.allow = async () => {
  const chapiFrame = await $('iframe');
  await chapiFrame.waitForExist();
  await expect(chapiFrame).toHaveAttrContaining('src', 'https://authn.io/mediator');
  await browser.switchToFrame(chapiFrame);
  const allowBtn = await $('button*=Allow');
  await allowBtn.waitForClickable();
  await allowBtn.click();
  await chapiFrame.waitForExist({ reverse: true });
  await browser.switchToFrame(null);
};
