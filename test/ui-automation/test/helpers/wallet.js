/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict';

const constants = require('./constants');
const { allow } = require('./chapi');

const DIDS = constants.dids;
let signedUpUserEmail;

/*************************** Public API ******************************/

exports.signUp = async ({ email }, isCHAPIEnabled = false) => {
  // login and consent
  await _getSignUp(email);

  if (isCHAPIEnabled) {
    // register chapi
    await allow();
  }

  // wait for default default vault to load successfully
  await _waitForDefaultVault();
  signedUpUserEmail = email;
};

exports.signIn = async (email = signedUpUserEmail) => {
  await _signIn(email);
};

exports.createOrbDID = async () => {
  await _createOrbDID();
};

exports.importDID = async ({ method }) => {
  await _importDID({ method: method });
};

exports.updatePreferences = async () => {
  await _updatePreferences();
};

exports.useJWTCredentials = async () => {
  await _useJWTCredentials();
};

exports.authenticate = async ({ did }) => {
  await _didAuth({ method: did });
};

exports.storeCredentials = async () => {
  await _acceptCredentials();
};

exports.presentCredentials = async () => {
  await _sendCredentials();
};

exports.addNewVault = async (vaultName) => {
  await _addNewVault(vaultName);
};

exports.renameVault = async (oldName, newName) => {
  await _renameVault(oldName, newName);
};

exports.removeVault = async (name) => {
  await _removeVault(name);
};

exports.vaultNameInput = async (vaultName) => {
  await _vaultNameInput(vaultName);
};

exports.createVault = async () => {
  await _createVault();
};

exports.cancelAddVault = async () => {
  await _cancelAddVault();
};

exports.validationError = async (msg) => {
  await _validationError(msg);
};

exports.validateVaultNameWithSpaces = async (actualVal, expectedVal) => {
  await _validateVaultNameWithSpaces(actualVal, expectedVal);
};

exports.validateUserInput = async (vaultName, errMsg) => {
  await _validateUserInput(vaultName, errMsg);
};

exports.waitForCredentials = async () => {
  await _waitForCredentials();
};

exports.didConnect = async () => {
  const didConnectBtn = await $('#didconnect');
  await didConnectBtn.waitForClickable();
  await didConnectBtn.click();

  const successMsg = await $('div*=CONGRATULATIONS ');
  await expect(successMsg).toExist();
};

exports.signOut = async () => {
  await _signOutWallet();
};

exports.checkStoredCredentials = async () => {
  await _checkStoredCredentials();
};

exports.changeLocale = async () => {
  await _changeLocale();
};

exports.validateCredentialDetails = async (vcData) => {
  for (const data of vcData) {
    // TODO need a better way to validate name and value matches rather than text existence on the screen
    const name = await $('td*=' + data.name);
    await expect(name).toExist();

    const val = await $('td*=' + data.value);
    await expect(val).toExist();
  }
};

exports.deleteCredential = async () => {
  const flyoutMenuImage = await $('#credential-details-flyout-button');
  await flyoutMenuImage.waitForClickable();
  await flyoutMenuImage.click();

  const deleteCredentialList = await $('#deleteCredential');
  await deleteCredentialList.waitForClickable();
  await deleteCredentialList.click();

  const deleteButton = await $('#delete-credential-button');
  await deleteButton.waitForClickable();
  await deleteButton.click();
};

exports.renameCredential = async (credentialName, errMsg) => {
  await _renameCredential(credentialName, errMsg);
};

async function _renameCredential(credentialName, errMsg) {
  const flyoutMenuImage = await $('#credential-details-flyout-button');
  await flyoutMenuImage.waitForClickable();
  await flyoutMenuImage.click();

  const renameCredentialList = await $('#renameCredential');
  await renameCredentialList.waitForClickable();
  await renameCredentialList.click();

  await _updateCredentialName(credentialName);
  const renameCredButton = await $('#rename-credential-button');
  await renameCredButton.waitForClickable();
  await renameCredButton.click();
  // validating rename credential name with the expected error Msg
  if (errMsg) {
    const errorMsg = await $('#input-CredentialName-error-msg');
    await errorMsg.waitForExist();
    await expect(errorMsg).toHaveText(errMsg);
    const dangerIcon = await $('.danger-icon');
    await expect(dangerIcon).toExist();
    await _cancelRenameCredential();
  }
}
/*************************** Helper functions ******************************/

async function _didAuth() {
  const authenticate = await $('#didauth');
  await authenticate.waitForClickable();
  await authenticate.click();
}

async function _acceptCredentials() {
  const storeBtn = await $('#storeVCBtn');
  await storeBtn.waitForClickable();
  await storeBtn.click();
}

async function _sendCredentials() {
  // share
  const shareBtn = await $('#share-credentials');
  await shareBtn.waitForClickable();
  await shareBtn.click();
}

async function _getSignUp(email) {
  const parentWindow = await browser.getWindowHandle();

  const signUpButton = await $('#mockbank');
  await signUpButton.waitForClickable();
  await signUpButton.click();

  await _getThirdPartyLogin(email, parentWindow);
}

async function _signOutWallet() {
  const signOutButton = await $('#signout-button');
  await signOutButton.waitForClickable();
  await signOutButton.click();
  await expect(browser).toHaveUrlContaining(browser.config.authURL + '/ui/sign-up');
}

async function _signIn(signedUpUserEmail) {
  const parentWindow = await browser.getWindowHandle();

  const signInButton = await $('#mockbank');
  await signInButton.waitForClickable();
  await signInButton.click();

  await _getThirdPartyLogin(signedUpUserEmail, parentWindow);
}

async function _changeLocale() {
  const localeSwitcherLink = await $('a*=Français');
  await localeSwitcherLink.waitForClickable();
  await localeSwitcherLink.click();

  const headingLink = await $('h1*=Inscrivez-vous. C’est gratuit!');
  await expect(headingLink).toExist();
  await expect(headingLink).toHaveText('Inscrivez-vous. C’est gratuit!');
}

async function _getThirdPartyLogin(email, parentWindow) {
  const windows = await browser.getWindowHandles();

  for (let i = 0; i < windows.length; i++) {
    if (windows[i] !== parentWindow) {
      await browser.switchToWindow(windows[i]);
      break;
    }
  }

  const emailInput = await $('#email');
  await emailInput.waitForExist();
  await emailInput.setValue(email);

  const oidcLoginButton = await $('#accept');
  await oidcLoginButton.waitForClickable();
  await oidcLoginButton.click();

  await browser.switchToWindow(parentWindow);
}

async function _waitForDefaultVault() {
  const defaultVault = await $('div*=Default Vault');
  await expect(defaultVault).toExist();
}

async function _waitForCredentials() {
  const credentialsLink = await $('#navbar-link-credentials');
  await credentialsLink.waitForClickable();
  await credentialsLink.click();
  const didResponse = await $('#loaded-credentials-container');
  await expect(didResponse).toExist();
  await expect(didResponse).toBeDisplayed();
}

async function _checkStoredCredentials() {
  const checkStoredCredential = await $('div*=Permanent Resident Card');
  await expect(checkStoredCredential).toExist();
}

async function _importDID({ method }) {
  const importDID = await $('label*=Import Any Digital Identity');
  await importDID.waitForClickable();
  await importDID.click();

  if (!DIDS[method]) {
    throw `couldn't find did method '${method} in test config'`;
  }

  const didInput = await $('#did-input');
  await didInput.waitForExist();
  await didInput.addValue(DIDS[method].did);

  const jwkType = await $('#JWK');
  await jwkType.waitForClickable();
  await jwkType.click();

  const privateKeyJWK = await $('#privateKeyStr');
  await privateKeyJWK.waitForExist();
  await privateKeyJWK.addValue(DIDS[method].pkjwk);

  const keyID = await $('#keyID');
  await keyID.waitForExist();
  await keyID.addValue(DIDS[method].keyID);

  const submit = await $('#saveDIDBtn');
  await submit.waitForClickable();
  await submit.click();

  const didResponse = await $('#save-anydid-success');
  await didResponse.waitForExist();
  await expect(didResponse).toHaveText('Saved your DID successfully.');
}

async function _createOrbDID() {
  const settingsTab = await $('a*=Settings');
  await settingsTab.waitForClickable();
  await settingsTab.click();

  const createOrbTab = await $('label*=Create ORB Digital Identity');
  await createOrbTab.waitForClickable();
  await createOrbTab.click();

  // select key Type
  const keyType = await $('#select-key');
  await keyType.waitForExist();
  await keyType.addValue(DIDS.orb.keyType);

  // select signature Type
  const signType = await $('#select-signature-suite');
  await signType.waitForExist();
  await signType.addValue(DIDS.orb.signatureType);

  const submit = await $('#createDIDBtn');
  await submit.waitForClickable();
  await submit.click();

  const didResponse = await $('#create-did-success');
  await didResponse.waitForExist();
  await expect(didResponse).toHaveText('Saved your DID successfully.');
}

async function _updatePreferences() {
  const preferences = await $('label*=Digital Identity Preference');
  await preferences.waitForClickable();
  await preferences.click();

  const jwkType = await $('label*=JsonWebSignature2020');
  await jwkType.waitForClickable();
  await jwkType.click();

  const submit = await $('button*=Update Preferences');
  await submit.waitForClickable();
  await submit.click();

  const successMessage = await $('#update-preferences-success');
  await expect(successMessage).toExist();
}

async function _useJWTCredentials() {
  const settingsTab = await $('a*=Settings');
  await settingsTab.waitForClickable();
  await settingsTab.click();

  const ldpType = await $('label*=JWT');
  await ldpType.waitForClickable();
  await ldpType.click();

  const submit = await $('button*=Update Preferences');
  await submit.waitForClickable();
  await submit.click();

  const successMessage = await $('#update-preferences-success');
  await expect(successMessage).toExist();
}

async function _addNewVault(vaultName) {
  // User clicks on Add Vault button
  const addVaultButton = await $('#add-new-vault-button');
  await addVaultButton.waitForClickable();
  await addVaultButton.click();
  await _vaultNameInput(vaultName);
  await _createVault();
}

async function _vaultNameInput(vaultName) {
  const addVaultInput = await $('#input-VaultName');
  await addVaultInput.waitForClickable();
  await addVaultInput.click();
  await addVaultInput.setValue(vaultName);
}

async function _updateCredentialName(credentialName) {
  const renameCredInput = await $('#input-CredentialName');
  await renameCredInput.waitForClickable();
  await renameCredInput.click();
  await renameCredInput.setValue(credentialName);
}

async function _createVault() {
  const addAction = await $('.btn-primary*=Add');
  await addAction.waitForClickable();
  await addAction.click();
}

async function _validationError(msg) {
  const errorMsg = await $('#input-VaultName-error-msg');
  await expect(errorMsg).toHaveText(msg);
  const dangerIcon = await $('.danger-icon');
  await expect(dangerIcon).toExist();
}

async function _cancelAddVault() {
  const cancelVaultButton = await $('.btn-outline*=Cancel');
  await cancelVaultButton.click();
  await expect(browser).toHaveUrlContaining('vaults');
}

async function _cancelRenameCredential() {
  const cancelCredentialButton = await $('.btn-outline*=Cancel');
  await cancelCredentialButton.click();
  await expect(browser).toHaveUrlContaining('credentials');
}

async function _validateUserInput(vaultName, errMsg) {
  await _addNewVault(vaultName);
  await _validationError(errMsg);
  await _cancelAddVault();
}

async function _validateVaultNameWithSpaces(actualVal, expectedVal) {
  await _addNewVault(actualVal);
  const vaultCard = await $(`#vault-card-${expectedVal.replaceAll(' ', '-')}`);
  await expect(vaultCard).toExist();
}

async function _renameVault(oldName, newName) {
  const vaultFlyoutButton = await $(`#vaults-flyout-menu-button-${oldName.replaceAll(' ', '-')}`);
  await vaultFlyoutButton.waitForClickable();
  await vaultFlyoutButton.click();

  const renameVaultButton = await $('#renameVault');
  await renameVaultButton.waitForClickable();
  await renameVaultButton.click();

  await _vaultNameInput(newName);

  const renameButton = await $('.btn-primary*=Rename');
  await renameButton.waitForClickable();
  await renameButton.click();
}

async function _removeVault(name) {
  const vaultFlyoutButton = await $(`#vaults-flyout-menu-button-${name.replaceAll(' ', '-')}`);
  await vaultFlyoutButton.waitForClickable();
  await vaultFlyoutButton.click();

  const renameVaultButton = await $('#delete-vault-flyout-button');
  await renameVaultButton.waitForClickable();
  await renameVaultButton.click();

  const renameButton = await $('.btn-danger*=Delete');
  await renameButton.waitForClickable();
  await renameButton.click();
}
