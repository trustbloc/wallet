/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import {wallet} from "../helpers";

describe("TrustBloc Wallet - Add/Rename/Delete Vault flow", () => {
  const ctx = {
    email: `ui-aut-${new Date().getTime()}@test.com`,
  };
  const vaultEndpoint =  'vaults';
  const validVaultName = 'Test Vault 1';
  const invalidVaultName ='Vault!@';

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
    it(`User successfully adds vault`, async function () {
        this.timeout(300000);
         await wallet.addNewVault();
         await wallet.vaultNameInput(validVaultName);
         await wallet.createVault();
         // Look for checkmark in the input to show user name is entered correctly
         const checkMarkIcon = await $('.checkmark-icon')
         await checkMarkIcon.waitForExist();

         // Show loading icon on the Add button.
         const loadingIcon = await $('.animate-spin');
         await loadingIcon.waitForExist();
         expect(browser.config.walletURL).toHaveValue(vaultEndpoint);
    });
  it(`Error handling of invalid user input`, async function () {
    this.timeout(300000);
    await wallet.validateUserInput("", "Can't be empty. Please enter a name." );
    await wallet.validateUserInput(invalidVaultName, "Must use letters (A-Z) and/or numbers (1-9)" );
    await wallet.validateUserInput(validVaultName, "There's already a vault with that name. Try something else." );
  });
  //TODO: Issue-1389 Add as many vault scenario
  //TODO: Issue-1393 Add locale support
});
