#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@all
@login
Feature: Login

  Scenario: User login
    When the user clicks on the Login button
     And the user is redirected to the OIDC provider
     And the user is authenticated
     And the user consents to sharing their identity data
    Then the user is redirected to the wallet's dashboard
     And the user can retrieve their profile

  Scenario: User logout
    When the user is logged in
     And the user logs out
    Then the user cannot access their profile
