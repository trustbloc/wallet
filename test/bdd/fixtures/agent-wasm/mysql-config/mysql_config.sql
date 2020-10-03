/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

\! echo "Configuring MySQL users...";

/*
Hydra
*/
CREATE USER 'hydra'@'%' IDENTIFIED BY 'hydra-pwd';
CREATE DATABASE hydra;
GRANT ALL PRIVILEGES ON hydra.* TO 'hydra'@'%';

