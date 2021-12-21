/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

"use strict";

const constants = require("./constants");
const { allow } = require("./chapi");

const DIDS = constants.dids;
const timeout = 60000;
let signedUpUserEmail;

/*************************** Public API ******************************/

exports.init = async ({ email }) => {
  // login and consent
  await _getSignUp(email);
  // register chapi
  await allow();

  // wait for credentials to load
  await _waitForCredentials();

  signedUpUserEmail = email;
};

exports.signIn = async () => {
  await _signIn(signedUpUserEmail);
};

exports.createOrbDID = async () => {
  await _createOrbDID();
};

exports.importDID = async ({ method, keyFormat }) => {
  await _importDID({ method: method });
};

exports.updatePreferences = async () => {
  await _updatePreferences();
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

exports.didConnect = async () => {
  const didConnectBtn = await $("#didconnect");
  await didConnectBtn.waitForExist();
  await didConnectBtn.waitForClickable();
  await didConnectBtn.click();

  const successMsg = await $("div*=CONGRATULATIONS ");
  await successMsg.waitForExist();
};

exports.logout = async () => {
  await _logoutWallet();
};

exports.checkStoredCredentials = async () => {
  await _checkStoredCredentials();
};

exports.changeLocale = async () => {
  await _changeLocale();
};

exports.validateCredentialDetails = async (vcData) => {
  for (const data of vcData) {
    // TODO need a better way to validate name and value matches rather than text existance on the screen
    const name = await $("td*=" + data.name);
    await name.waitForExist();

    const val = await $("td*=" + data.value);
    await val.waitForExist();
  }
};

exports.deleteCredential = async () => {
  const flyoutMenuImage = await $("#credential-details-flyout-button");
  await flyoutMenuImage.waitForExist();
  await flyoutMenuImage.waitForClickable();
  await flyoutMenuImage.click();

  const deleteCredentialList = await $("#deleteCredential");
  await deleteCredentialList.waitForExist();
  await deleteCredentialList.waitForClickable();
  await deleteCredentialList.click();

  const deleteButton = await $("#delete-credential-button");
  await deleteButton.waitForExist();
  await deleteButton.waitForClickable();
  await deleteButton.click();
};

/*************************** Helper functions ******************************/

async function _didAuth({ method = "trustbloc" } = {}) {
  const authenticate = await $("#didauth");
  await authenticate.waitForExist();
  await authenticate.click();
}

async function _acceptCredentials() {
  const storeBtn = await $("#storeVCBtn");
  await storeBtn.waitForExist();
  await storeBtn.waitForClickable();
  await storeBtn.click();
}

async function _sendCredentials() {
  // share
  const shareBtn = await $("#share-credentials");
  await shareBtn.waitForExist();
  await shareBtn.waitForClickable();
  await shareBtn.click();
}

async function _getSignUp(email) {
  const signUpButton = await $("#mockbank");
  await signUpButton.waitForExist();
  await signUpButton.click();
  await _getThirdPartyLogin(email);
}

async function _logoutWallet() {
  const logOutButton = await $("button*=Sign Out");
  await logOutButton.waitForExist();
  await logOutButton.click();

  // wait for logout to complele and go to signup page
  await browser.waitUntil(async () => {
    const headingLink = await $("h1*=Sign up.");
    expect(headingLink).toHaveValue("Sign up.");
    return true;
  });
}

async function _signIn(signedUpUserEmail) {
  const signInLink = await $("a*=Sign in");
  await signInLink.waitForExist();
  await signInLink.click();
  await browser.waitUntil(async () => {
    const signInButton = await $("#mockbank");
    await signInButton.waitForExist();
    await signInButton.click();
    await _getThirdPartyLogin(signedUpUserEmail);
    return true;
  });
  return true;
}

async function _changeLocale() {
  const localeSwitcherLink = await $("a*=Français");
  await localeSwitcherLink.waitForExist();
  await localeSwitcherLink.click();
  await browser.waitUntil(async () => {
    const headingLink = await $("h1*=Inscrivez-vous. C’est gratuit!");
    expect(headingLink).toHaveValue("Inscrivez-vous. C’est gratuit!");
    return true;
  });
}

async function _getThirdPartyLogin(email) {
  await browser.waitUntil(async () => {
    try {
      await browser.switchWindow("Login Page");
    } catch (err) {
      console.warn("[warn] switch window to login page : ", err.message);
      return false;
    }
    return true;
  });

  await browser.waitUntil(async () => {
    let emailInput = await $("#email");
    await emailInput.waitForExist();
    expect(emailInput).toHaveValue("john.smith@example.com");
    await emailInput.setValue(email);
    return true;
  });

  const oidcLoginButton = await $("#accept");
  await oidcLoginButton.click();

  await browser.switchWindow(browser.config.walletURL);
  await browser.waitUntil(async () => {
    let title = await $("iframe");
    await title.waitForExist({ timeout, interval: 5000 });
    return true;
  });
}

async function _waitForCredentials() {
  await browser.waitUntil(async () => {
    const credentialsLink = await $("#navbar-link-credentials");
    await credentialsLink.click();
    let didResponse = await $("#loaded-credentials-container");
    await didResponse.waitForExist({ timeout, interval: 5000 });
    expect(didResponse).toBeDisplayed();
    return true;
  });
}

async function _checkStoredCredentials() {
  const checkStoredCredential = await $("div*=Permanent Resident Card");
  await checkStoredCredential.waitForExist();
  return true;
}

async function _importDID({ method }) {
  const settingsTab = await $("a*=Settings");
  await settingsTab.waitForExist();
  await settingsTab.click();

  const importDID = await $("label*=Import Any Digital Identity");
  await importDID.waitForExist();
  await importDID.click();

  if (!DIDS[method]) {
    throw `couldn't find did method '${did} in test config'`;
  }

  const didInput = await $("#did-input");
  await didInput.addValue(DIDS[method].did);

  const jwkType = await $("#JWK");
  await jwkType.click();

  const privateKeyJWK = await $("#privateKeyStr");
  await privateKeyJWK.addValue(DIDS[method].pkjwk);

  const keyID = await $("#keyID");
  await keyID.addValue(DIDS[method].keyID);

  const submit = await $("#saveDIDBtn");
  await submit.click();

  await browser.waitUntil(async () => {
    let didResponse = await $("#save-anydid-success");
    await didResponse.waitForExist({ timeout, interval: 2000 });
    expect(didResponse).toHaveText("Saved your DID successfully.");
    return true;
  });
}

async function _createOrbDID() {
  const settingsTab = await $("a*=Settings");
  await settingsTab.waitForExist();
  await settingsTab.click();

  const createOrbTab = await $("label*=Create ORB Digital Identity");
  await createOrbTab.waitForClickable();
  await createOrbTab.click();

  // select key Type
  const keyType = await $("#select-key");
  await keyType.addValue(DIDS.orb.keyType);

  // select signature Type
  const signType = await $("#select-signature-suite");
  await signType.addValue(DIDS.orb.signatureType);

  const submit = await $("#createDIDBtn");
  await submit.click();

  await browser.waitUntil(async () => {
    let didResponse = await $("#create-did-success");
    await didResponse.waitForExist({ timeout, interval: 2000 });
    expect(didResponse).toHaveText("Saved your DID successfully.");
    return true;
  });
}

async function _updatePreferences() {
  const settingsTab = await $("a*=Settings");
  await settingsTab.waitForExist();
  await settingsTab.click();

  const preferences = await $("label*=Digital Identity Preference");
  await preferences.waitForExist();
  await preferences.click();

  const jwkType = await $("label*=JsonWebSignature2020");
  await jwkType.click();

  const submit = await $("button*=Update Preferences");
  await submit.click();

  // TODO validate success message
}
