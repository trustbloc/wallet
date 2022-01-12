/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import {wallet} from "../helpers";

const vault = new Map();
vault.set("Test Vault 1", {
    "@context":["https://w3id.org/wallet/v1",
        "https://trustbloc.github.io/context/wallet/collections-v1.jsonld"],
    "id":"1",
    "name":"Test Vault 1",
    "type":"Vault"
});
vault.set("Test Vault 2", {
    "@context":["https://w3id.org/wallet/v1",
        "https://trustbloc.github.io/context/wallet/collections-v1.jsonld"],
    "id":"2",
    "name":"Test Vault 2",
    "type":"Vault"
});
vault.set("Test Vault 3", {
    "@context":["https://w3id.org/wallet/v1",
        "https://trustbloc.github.io/context/wallet/collections-v1.jsonld"],
    "id":"3",
    "name":"Test Vault 3",
    "type":"Vault"
});
vault.set("Test Vault 4", {
    "@context":["https://w3id.org/wallet/v1",
        "https://trustbloc.github.io/context/wallet/collections-v1.jsonld"],
    "id":"4",
    "name":"Test Vault 4",
    "type":"Vault"
});

describe("TrustBloc Wallet - Add/Rename/Delete Vault flow", () => {
    const ctx = {
        email: `ui-aut-${new Date().getTime()}@test.com`,
    };
    const validVaultName = "Test Vault 4";
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
    it(`User successfully add vaults`, async function () {
        this.timeout(600000);
        for (const [key, value] of vault.entries()) {
        await wallet.addNewVault(key);
        // Look for checkmark in the input to show user name is entered correctly
        const checkMarkIcon = await $('.checkmark-icon')
        await checkMarkIcon.waitForExist();

        await browser.pause(6000);
        const vaultCard = await $('.vaultContainer');
        await vaultCard.waitForExist();
        expect(vaultCard).toHaveTextContaining(key);
        await browser.refresh();
        }
    });
    it(`Error handling of invalid user input`, async function () {
        await wallet.validateUserInput("", "Can't be empty. Please enter a name." );
        await wallet.validateUserInput(invalidVaultName, "Must use letters (A-Z) and/or numbers (1-9)" );
        await wallet.validateUserInput(validVaultName, "There's already a vault with that name. Try something else." );
    });
    it(`User enter vault name with spaces`, async function () {
        // vault name with multiple trailing spaces
        await wallet.validateVaultNameWithSpaces('Testing Vault with trailing spaces   ', 'Testing Vault with trailing spaces');
        // vault name with repetitive spaces
        await wallet.validateVaultNameWithSpaces('Testing Vault   with repetitive spaces', 'Testing Vault with repetitive spaces');
        // vault name with leading spaces
        await wallet.validateVaultNameWithSpaces('   Testing Vault with leading spaces', 'Testing Vault with leading spaces');
        // vault name with all the spaces combined
        await wallet.validateVaultNameWithSpaces('   Testing    vault  with  all   spaces   ', 'Testing vault with all spaces');
    });
    it(`User changes vault locale (${ctx.email})`, async function () {
        this.timeout(90000);

        // 1. Navigate to Wallet Website
        await browser.navigateTo(browser.config.walletURL);

        // 2. Change locale
        const localeSwitcherLink = await $("a*=Français");
        await localeSwitcherLink.waitForExist();
        await localeSwitcherLink.click();
        await browser.waitUntil(async () => {
            // Verifying lang attribute is set to French on the page
            await expect(browser).toHaveUrlContaining(browser.config.walletURLFrench);
            const headingLink = await $("h3*=Chambres fortes");
            expect(headingLink).toHaveValue("Chambres fortes");
            const addVaultButton = await $("#add-new-vault-button");
            expect(addVaultButton).toHaveValue("Ajouter une chambre forte");
            const allVaultsButton = await $("#all-vaults-button");
            expect(allVaultsButton).toHaveValue("Toutes les chambres fortes");
            return true;
        });
    });

});