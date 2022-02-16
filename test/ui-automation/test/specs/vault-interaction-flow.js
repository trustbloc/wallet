/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

import { wallet } from "../helpers";

const vault = new Map();
vault.set("Test Vault 1", {
  "@context": [
    "https://w3id.org/wallet/v1",
    "https://trustbloc.github.io/context/wallet/collections-v1.jsonld",
  ],
  id: "1",
  name: "Test Vault 1",
  type: "Vault",
});
vault.set("Test Vault 2", {
  "@context": [
    "https://w3id.org/wallet/v1",
    "https://trustbloc.github.io/context/wallet/collections-v1.jsonld",
  ],
  id: "2",
  name: "Test Vault 2",
  type: "Vault",
});
vault.set("Test Vault 3", {
  "@context": [
    "https://w3id.org/wallet/v1",
    "https://trustbloc.github.io/context/wallet/collections-v1.jsonld",
  ],
  id: "3",
  name: "Test Vault 3",
  type: "Vault",
});
vault.set("Test Vault 4", {
  "@context": [
    "https://w3id.org/wallet/v1",
    "https://trustbloc.github.io/context/wallet/collections-v1.jsonld",
  ],
  id: "4",
  name: "Test Vault 4",
  type: "Vault",
});

describe("TrustBloc Wallet - Add/Rename/Delete Vault flow", () => {
  const ctx = {
    email: `ui-aut-${new Date().getTime()}@test.com`,
  };
  const validVaultName = "Default Vault";
  const invalidVaultName = "Vault!@";
  // runs once before the first test in this block
  before(async () => {
    await browser.reloadSession();
    await browser.maximizeWindow();
  });

  it(`User signs up (${ctx.email})`, async function () {
    // 1. Navigate to Wallet Website
    await browser.navigateTo(browser.config.walletURL);

    // 2. Initialize Wallet (register/sign-up/etc.)
    await wallet.signUp(ctx);
  });
  it(`User dismisses welcome message banner`, async function () {
    const welcomeMessageCloseButton = await $("#welcome-banner-close-button");
    await welcomeMessageCloseButton.waitForExist();
    await welcomeMessageCloseButton.click();
    const dashboardContent = await $("#dashboard-content");
    await dashboardContent.waitForExist();
    await expect(dashboardContent).toHaveChildren(1);
  });
  it(`User successfully adds vaults`, async function () {
    for (const [key, value] of vault.entries()) {
      await wallet.addNewVault(key);
      const vaultCard = await $(
        `#vault-card-${value.name.replaceAll(" ", "-")}`
      );
      await vaultCard.waitForExist();
    }
  });
  it("User successfully renames vaults", async function () {
    for (const [key, value] of vault.entries()) {
      await wallet.renameVault(key, `${key} renamed`);
      const vaultCard = await $(
        `#vault-card-${value.name.replaceAll(" ", "-")}-renamed`
      );
      await vaultCard.waitForExist();
    }
  });
  it("User successfully removes vaults", async function () {
    // Number of vault cards
    let numOfVaults = vault.entries().length + 3;
    for (const [key, value] of vault.entries()) {
      await wallet.removeVault(`${key} renamed`);
      const vaultsLoadedContainer = await $("#vaults-loaded");
      await vaultsLoadedContainer.waitForExist();
      numOfVaults -= 1;
      await expect(vaultsLoadedContainer).toHaveChildren(numOfVaults);
    }
  });
  it(`User enters invalid vault name`, async function () {
    await wallet.validateUserInput("", "Can't be empty. Please enter a name.");
    await wallet.validateUserInput(
      invalidVaultName,
      "Must use letters (A-Z) and/or numbers (1-9)"
    );
    await wallet.validateUserInput(
      validVaultName,
      "There's already a vault with that name. Try something else."
    );
  });
  it(`User enters vault name with spaces`, async function () {
    // vault name with multiple trailing spaces
    await wallet.validateVaultNameWithSpaces(
      "Testing Vault with trailing spaces   ",
      "Testing Vault with trailing spaces"
    );
    // vault name with repetitive spaces
    await wallet.validateVaultNameWithSpaces(
      "Testing Vault   with repetitive spaces",
      "Testing Vault with repetitive spaces"
    );
    // vault name with leading spaces
    await wallet.validateVaultNameWithSpaces(
      "   Testing Vault with leading spaces",
      "Testing Vault with leading spaces"
    );
    // vault name with all the spaces combined
    await wallet.validateVaultNameWithSpaces(
      "   Testing    vault  with  all   spaces   ",
      "Testing vault with all spaces"
    );
  });
  it(`User changes vault locale (${ctx.email})`, async function () {
    const localeSwitcherLink = await $("a*=Français");
    await localeSwitcherLink.waitForExist();
    await localeSwitcherLink.click();
    await expect(browser).toHaveUrlContaining(browser.config.walletURLFrench);
    const headingLink = await $("h3*=Chambres fortes");
    await expect(headingLink).toExist();
    const addVaultCard = await $("span*=Ajouter une chambre forte");
    await expect(addVaultCard).toExist();
    const allVaultsButton = await $("span*=Toutes les chambres fortes");
    await expect(allVaultsButton).toExist();
  });
  it("User validates dismissed welcome banner", async function () {
    await wallet.signOut(ctx);
    await wallet.signIn(ctx.email);
    const dashboardContent = await $("#dashboard-content");
    await dashboardContent.waitForExist();
    await expect(dashboardContent).toHaveChildren(1);
  });
});
