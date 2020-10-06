#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@all
@login
Feature: Login

  Scenario: New User
    When the user clicks on the Login button
     And the user is redirected to the OIDC provider
     And the user is authenticated
     And the user consents to sharing their identity data
    Then the user is redirected to the wallet's dashboard
