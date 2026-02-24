# DOCSIS TLV Reference

## 1. Overview

This document catalogs all DOCSIS provisioning TLVs from the CableLabs Assigned Names and Numbers registry (CL-SP-CANN-I24). Canonical names are the CANN-registered names that will be used as JSON schema property names in OpenDCI.

**Sources:**

- **CANN** -- CL-SP-CANN-I24-250320, Section 11.1 (DOCSIS Provisioning TLV Number Assignment Registry)
- **MULPIv4.0** -- CM-SP-MULPIv4.0-I11-260219, Annex C (Configuration File and Registration TLV Encodings)
- Individual extension specifications: L2VPN, eRouter, eDOCSIS, PacketCable, DSG, TEI, DPoE, DOCSIS SYNC

TLVs are used in CM configuration files and MAC Management messages. They are also used in the RPHY GCP Protocol. The Type-Length-Value encoding allows provisioning systems to configure cable modems and CMTS equipment with network parameters, QoS settings, security credentials, and service flow definitions.

---

## 2. Top-Level DOCSIS TLVs (CANN 11.1)

This table lists all top-level TLV types registered in CANN Section 11.1, cross-referenced with MULPIv4.0 Annex C Table 117 for length, config file applicability, and section references.

### Types 0--106

| Type | CANN Name | Data Type | Length | Cfg File | Spec | Reference |
|------|-----------|-----------|--------|----------|------|-----------|
| 0 | Pad | - | - | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.2 |
| 1 | Downstream Frequency | uint32 | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.1 |
| 2 | Upstream Channel ID | uint8 | 1 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.2 |
| 3 | Network Access Control Object | uint8 | 1 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.3 |
| 4 | DOCSIS 1.0 Class of Service | - | - | - | DOCSIS 1.0 · MULPIv4.0 I11 | *(deprecated)* |
| 5 | Modem Capabilities Encoding | encapsulated | n | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.1 |
| 6 | CM Message Integrity Check (MIC) | hexstring | 16 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.5 |
| 7 | CMTS Message Integrity Check (MIC) | hexstring | 16 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.6 |
| 8 | Vendor ID Encoding | hexstring | 3 | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.2 |
| 9 | SW Upgrade Filename | string | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.3 |
| 10 | SNMP Write Access Control | hexstring | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.4 |
| 11 | SNMP MIB Object | encapsulated | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.5 |
| 12 | Modem IP Address | ipv4Address | 4 | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.3 |
| 13 | Service(s) Not Available Response | hexstring | 3 | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.4 |
| 14 | CPE Ethernet MAC Address | macAddress | 6 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.6 |
| 15 | Telephone Settings Option | - | - | - | DOCSIS 1.0 · MULPIv4.0 I11 | *(deprecated)* |
| 16 | *(unassigned)* | - | - | - | - | - |
| 17 | Baseline Privacy (Security) | hexstring | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.3.1.1 |
| 18 | Max Number of CPEs | uint8 | 1 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.7 |
| 19 | TFTP Server Timestamp | uint32 | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.8 |
| 20 | TFTP Server Provisioned Modem IPv4 Address | ipv4Address | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.9 |
| 21 | SW Upgrade IPv4 TFTP Server | ipv4Address | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.7 |
| 22 | Upstream Packet Classification | encapsulated | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.1.1 |
| 23 | Downstream Packet Classification | encapsulated | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.1.3 |
| 24 | Upstream Service Flow | encapsulated | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.2.1 |
| 25 | Downstream Service Flow | encapsulated | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.2.2 |
| 26 | Payload Header Suppression | encapsulated | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.3 |
| 27 | HMAC-Digest | hexstring | 20 | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.4.1 |
| 28 | Maximum Number of Classifiers | uint16 | 2 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.16 |
| 29 | Privacy Enable | uint8 | 1 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.17 |
| 30 | Authorization Block | hexstring | n | No | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.4.2 |
| 31 | Key Sequence Number | uint8 | 1 | No | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.4.3 |
| 32 | Manufacturer Code Verification Certificate | hexstring | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.10 |
| 33 | Co-Signer Code Verification Certificate | hexstring | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.11 |
| 34 | SNMPv3 Kickstart Value | hexstring | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.9 |
| 35 | Subscriber Mgmt Control | hexstring | 3 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.19.1 |
| 36 | Subscriber Mgmt CPE IPv4 List | hexstring | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.19.2 |
| 37 | Subscriber Mgmt Filter Groups | hexstring | 8 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.19.4 |
| 38 | SNMPv3 Notification Receiver | hexstring | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.12 |
| 39 | Enable 2.0 Mode | uint8 | 1 | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.20 |
| 40 | Enable Test Modes | uint8 | 1 | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.20 |
| 41 | Downstream Channel List | hexstring | n | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.22 |
| 42 | Static Multicast MAC Address | macAddress | 6 | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.23 |
| 43 | DOCSIS Extension Field | encapsulated | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.18 |
| 44 | Vendor Specific Capabilities | hexstring | n | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.5 |
| 45 | Downstream Unencrypted Traffic (DUT) Filtering | encapsulated | n | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.24 |
| 46 | Transmit Channel Configuration (TCC) | encapsulated | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.1 |
| 47 | Service Flow SID Cluster Assignment | encapsulated | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.2 |
| 48 | Receive Channel Profile | encapsulated | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.3.1 |
| 49 | Receive Channel Configuration | encapsulated | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.3.1 |
| 50 | DSID Encodings | encapsulated | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.3.9 |
| 51 | Security Association Encoding | encapsulated | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.5 |
| 52 | Initializing Channel Timeout | uint16 | 2 | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.6 |
| 53 | SNMPv1v2c Coexistence | encapsulated | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.13 |
| 54 | SNMPv3 Access View Configuration | encapsulated | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.14 |
| 55 | SNMP CPE Access Control | uint8 | 1 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.15 |
| 56 | Channel Assignment Configuration Settings | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.25 |
| 57 | CM Initialization Reason | uint8 | 1 | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.3.6 |
| 58 | SW Upgrade IPv6 TFTP Server | ipv6Address | 16 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.8 |
| 59 | TFTP Server Provisioned Modem IPv6 Address | ipv6Address | 16 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.10 |
| 60 | Upstream Drop Packet Classification | encapsulated | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.2.1.2 |
| 61 | Subscriber Mgmt CPE IPv6 Prefix List | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.19.3 |
| 62 | Upstream Drop Classifier Group ID | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.26 |
| 63 | Subscriber Mgmt Control Max CPE IPv6 Prefix | uint16 | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.19.5 |
| 64 | CMTS Static Multicast Session Encoding | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.27 |
| 65 | L2VPN MAC Aging Encoding | encapsulated | n | Yes | L2VPN I17 | C.1.1.28 |
| 66 | Management Event Control Encoding | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.16 |
| 67 | Subscriber Mgmt CPE IPv6 List | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.19.6 |
| 68 | Default Upstream Target Buffer Configuration | uint16 | 2 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.17 |
| 69 | MAC Address Learning Control Encoding | uint8 | 1 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.18 |
| 70 | Upstream Aggregate Service Flow Encodings | encapsulated | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.3 |
| 71 | Downstream Aggregate Service Flow Encodings | encapsulated | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.4 |
| 72 | Metro Ethernet Service Profile | encapsulated | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 / DPoE-MULPIv2.0 I14 | C.2.2.12 |
| 73 | Network Timing Profile | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 / DPoE-MULPIv2.0 I14 | C.1.2.19 |
| 74 | Energy Management Parameter Encoding | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.30 |
| 75 | Energy Mgt. Mode Indicator | uint8 | 1 | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.4.4 |
| 76 | CM Upstream AQM Disable | uint8 | 1 | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.20 |
| 77 | DOCSIS Time Protocol Encoding | encapsulated | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.6 |
| 78 | Energy Management Identifier List for CM | hexstring | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.1.30.4 |
| 79 | UNI Control Encoding | hexstring | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.3.3 |
| 80 | Energy Management -- DOCSIS Light Sleep Encodings | encapsulated | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.4.5 |
| 81 | Manufacturer CVC Chain | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.21 |
| 82 | Co-signer CVC Chain | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.22 |
| 83 | L2CP Management | encapsulated | 1 | Yes | DPoE-MULPIv2.0 I14 | C.1.1.31 |
| 84 | Diplexer Band Edge | hexstring | 9 | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.23 |
| 85 | FDX Transmission Group Assignment | encapsulated | n | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.4.6 |
| 86 | FDX Reset | uint8 | 1 | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.4.7 |
| 87 | CM Echo Cancellation Training Control | encapsulated | n | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.7 |
| 88 | QoS Framework for DOCSIS Encodings | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.8 |
| 89 | Extended SID Cluster Assignment | encapsulated | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.3.1.70 |
| 90 | Primary Service Flow Indicator | encapsulated | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.3.7 |
| 91 | Low Latency Disable | uint8 | 1 | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.1.32 |
| 92 | Distributed HQoS Enable | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.1.33 |
| 93 | Upstream Enhanced HQoS ASF | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.5 |
| 94 | Downstream Enhanced HQoS ASF | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.6 |
| 95 | DHQoS ASF SID Bundle Assignment | hexstring | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.5.8 |
| 96 | Advanced Diplexer Band Edge | hexstring | n | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.2.24 |
| 97 | Advanced Band Plan Support | uint8 | 1 | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.2.25 |
| 98 | DOCSIS Sync Capabilities | encapsulated | n | No | SYNC I03 | [SYNC] |
| 99 | DOCSIS CM System Information | encapsulated | n | No | SYNC I03 | [SYNC] |
| 100 | Sync DSID Assignment | encapsulated | n | No | SYNC I03 | [SYNC] |
| 101 | DOCSIS Sync Configurations | hexstring | n | Yes | SYNC I03 | [SYNC] |
| 102 | PTP Address Configurations | hexstring | n | Yes | SYNC I03 | [SYNC] |
| 103 | CM SSH Server Configuration Settings | hexstring | n | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.3.1.2 |
| 104 | Security Configuration Settings | hexstring | n | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.3.1.3 |
| 105 | Extended Modem Capabilities | encapsulated | n | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.3.1 |
| 106 | FDX Downstream Upper Band Edge | uint16 | 2 | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.2.26 |

> **Note:** Types 16 and 107--200 are reserved/unassigned.

### eSAFE Types (201--231)

| Type | CANN Name | Data Type | Length | Cfg File | Spec | Reference |
|------|-----------|-----------|--------|----------|------|-----------|
| 201 | ePS | hexstring | n | Yes | eDOCSIS I31 | [eDOCSIS] |
| 202 | eRouter | encapsulated | n | Yes | eRouter I22 | [eRouter] |
| 203--215 | *(Reserved)* | - | - | - | - | - |
| 216 | eMTA | hexstring | n | Yes | PacketCable PROV1.5 C01 | [PacketCable] |
| 217 | eSTB | hexstring | n | Yes | DSG I25 | [DSG] |
| 218 | *(Reserved)* | - | - | - | - | - |
| 219 | eTEA | encapsulated | n | Yes | TEI I06 | [TEI] |
| 220 | eDVA | hexstring | n | Yes | PacketCable RST-E-DVA C01 | [PacketCable 2.0] |
| 221 | eSG | hexstring | n | Yes | eDOCSIS I31 | [eDOCSIS] |
| 222--231 | *(Reserved)* | - | - | - | - | - |

> **Note:** Types 232--254 are reserved/unassigned.

### Type 255

| Type | CANN Name | Data Type | Length | Cfg File | Spec | Reference |
|------|-----------|-----------|--------|----------|------|-----------|
| 255 | End-of-Data | uint8 | - | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.1 |

---

## 3. Sub-TLV Sections

### 3.1 TLV 5 -- Modem Capabilities Encoding Sub-TLVs

CANN Section 11.1.1. These sub-TLVs are carried within TLV 5 (Modem Capabilities Encoding) and are sent by the CM to the CMTS during registration to advertise device capabilities.

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 5.1 | Concatenation Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.2 | DOCSIS Version | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.3 | Fragmentation Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.4 | Payload Header Suppression Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.5 | IGMP Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.6 | Privacy Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.7 | Downstream SAID Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.8 | Upstream Service Flow Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.9 | Optional Filtering Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.10 | Transmit Pre-Equalizer Taps | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.11 | Number of Transmit Pre-Equalizer Taps | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.12 | DCC Support | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | - |
| 5.13 | IP Filters Support | uint8 | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 5.14 | LLC Filters Support | uint8 | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 5.15 | Expanded Unicast SID Space | uint8 | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 5.16 | Ranging Hold-off Support | hexstring | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 5.17 | L2VPN Capability | uint8 | L2VPN I17 | - |
| 5.18 | L2VPN eSAFE Host Capability | hexstring | L2VPN I17 | - |
| 5.19 | Downstream Unencrypted Traffic (DUT) Filtering | uint8 | L2VPN I17 | - |
| 5.20 | Upstream Frequency Range Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.21 | Upstream SC-QAM Symbol Rate Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.22 | Selectable Active Code Mode 2 Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.23 | Code Hopping Mode 2 Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.24 | Multiple Transmit SC-QAM Channel Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.25 | 5.12 Msps UpstreamTransmit SC-QAM Channel Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.26 | 2.56 Msps Upstream Transmit SC-QAM Channel Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.27 | Total SID Cluster Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.28 | SID Clusters per Service Flow Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.29 | Multiple Receive SC-QAM Channel Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.30 | Total Downstream Service ID (DSID) Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.31 | Resequencing Downstream Service ID (DSID) Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.32 | Multicast Downstream Service ID (DSID) Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.33 | Multicast DSID Forwarding | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.34 | Frame Control Type Forwarding Capability | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.35 | DPV Capability | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.36 | Unsolicited Grant Service/Upstream Service Flow Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.37 | MAP and UCD Receipt Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.38 | Upstream Drop Classifier Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.39 | IPv6 Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.40 | Extended Upstream Transmit Power Capability | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.41 | Optional 802.1ad, 802.1ah, MPLS Classification Support | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.42 | D-ONU Capabilities | encapsulated | DPoE-MULPIv1.0 C01 | - |
| 5.42.1 | DPoE/G Version Number | uint8 | DPoE-MULPIv1.0 C01 | - |
| 5.42.2 | Number of Unicast LLIDs | uint8 | DPoE-MULPIv1.0 C01 | - |
| 5.42.3 | Number of Multicast LLIDs | uint8 | DPoE-MULPIv2.0 I14 | - |
| 5.42.4 | MESP Support | uint8 | DPoE-MULPIv2.0 I14 | - |
| 5.42.5 | Number of D-ONU Ports | uint8 | DPoE-MULPIv2.0 I14 | - |
| 5.42.6 | PON Data Rate Support | uint8 | DPoE-MULPIv2.0 I14 | - |
| 5.42.7 | Service OAM | uint8 | DPoE-MULPIv2.0 I14 | - |
| 5.42.10 | Number of T-CONTs Supported | uint8 | DPoG 1.0 | - |
| 5.42.11 | Total Number of (X)GEM Ports Supported | uint8 | DPoG 1.0 | - |
| 5.43 | Reserved | - | - | - |
| 5.44 | Energy Management Capabilities | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.45 | C-DOCSIS Capability Encoding | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.46 | CM-STATUS-ACK | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 5.47 | Energy Management Preferences | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.48 | Extended Packet Length Support Capability | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.49 | Multiple Receive OFDM Channel Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.50 | Multiple Transmit OFDMA Channel Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.51 | Downstream OFDM Profile Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.52 | Downstream OFDM Channel Subcarrier QAM Modulation Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.53 | Upstream OFDMA Channel Subcarrier QAM Modulation Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.54 | Downstream Lower Band Edge Configuration | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.55 | Downstream Upper Band Edge Configuration | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.56 | Diplexer Upstream Upper Band Edge Configuration | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.57 | DOCSIS Time Protocol Mode | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.58 | DOCSIS Time Protocol Performance Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.59 | Pmax | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.60 | Diplexer Downstream Lower Band Edge Options | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.61 | Diplexer Downstream Upper Band Edge Options | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.62 | Diplexer Upstream Upper Band Edge Options | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.63 | Advanced Band Plan Capability | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.64 | FDX DS State Lock -- Deprecated | - | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.65 | FDX Switching Software Timing Uncertainty | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.66 | FDX DS to US Switching Time | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.67 | FDX US to DS Switching Time | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.68 | -- | - | - | - |
| 5.69 | CWT RxMER Measurement Convergence Time | uint8 | - | - |
| 5.70 | -- | - | - | - |
| 5.71 | -- | - | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.72 | t-ds-reacquisition capability | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.73 | CWT Simultaneous Data Transmission Capability | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.74 | Extended Service Flow SID Cluster Assignments Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.75 | Echo Cancelling RBA Sub-band Direction Sets Supported | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.76 | Low Latency Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.77 | Absolute Queue-Depth Request Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.78 | Distributed HQoS Support | uint8 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 5.79 | Advanced Downstream Lower Band Edge Configuration | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.80 | Advanced Downstream Upper Band Edge Configuration | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.81 | Advanced Diplexer Upstream Upper Band Edge Configuration | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.82 | Advanced Diplexer Downstream Lower Band Edge Options List | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.83 | Advanced Diplexer Downstream Upper Band Edge Options List | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.84 | Advanced Diplexer Upstream Upper Band Edge Options List | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |
| 5.85 | Extended Power Options | uint8 | DOCSIS 4.0 · MULPIv4.0 I11 | - |

### 3.2 TLV 43 -- DOCSIS Extension Field Sub-TLVs (General Extension)

CANN Section 11.1.2. When TLV 43 carries OUI = 0xFFFFFF (General Extension Information), the following sub-TLVs are defined. See also [Section 4](#4-tlv-43----vendor-specific-information) for the full TLV 43 structure.

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 43.1 | CM Load Balancing Policy ID | uint32 | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.18.1.1 |
| 43.2 | CM Load Balancing Priority | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.18.1.2 |
| 43.3 | CM Load Balancing Group ID | uint32 | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.18.1.3 |
| 43.4 | CM Ranging Class ID Extension | uint16 | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.18.1.4 |
| 43.5 | L2VPN Encoding | encapsulated | L2VPN I17 | C.1.1.18.1.5 |
| 43.6 | Extended CMTS MIC Configuration Setting | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.6 |
| 43.7 | Source Address Verification (SAV) | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.7 |
| 43.8 | Cable Modem Attribute Masks | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.8 |
| 43.9 | IP Multicast Join Authorization | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.9 |
| 43.10 | Service Type Identifier | string | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.10 |
| 43.12 | DEMARC Auto-Configuration (DAC) | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.11 |

> **Note:** Sub-type 43.11 is unassigned.

### 3.3 TLV 43.5 -- L2VPN Encoding Sub-TLVs

CANN Section 11.1.2.1. These sub-TLVs are carried within TLV 43.5 (L2VPN Encoding) inside the General Extension Information field.

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 43.5.1 | VPN Identifier | hexstring | L2VPN I17 | - |
| 43.5.2 | NSI encapsulation format | encapsulated | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.1 | Other | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.2 | IEEE 802.1Q | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.3 | IEEE 802.1ad | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.4 | MPLS PW | encapsulated | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.4.1 | MPLS Pseudowire ID | uint32 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.4.2 | MPLS Peer IP address | ipv4Address | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.4.3 | Pseudowire Type | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.4.4 | MPLS Backup Pseudowire ID | uint32 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.4.5 | MPLS Backup Peer IP address | ipv4Address | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.5 | L2TPv3 Peer | ipv4Address | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.6 | IEEE 802.1ah Encapsulation | encapsulated | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.6.1 | IEEE 802.1ah Backbone Service Instance Tag (I-Tag) TCI | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.6.2 | IEEE 802.1ah Destination Backbone Edge Bridge (BEB) MAC Address (B-DA) | macAddress | DPoE-MULPIv1.0 C01 | - |
| 43.5.2.6.3 | 16-bit value of [802.1ah] B-Tag TCI | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.4 | 16-bit value of [802.1ah] I-Tag TPID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.5 | 3 bit I-PCP | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.6 | 1 bit I-DEI | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.7 | 1 bit I-UCA | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.8 | 24-bit value of [802.1ah] I-SID Backbone Service Instance Identifier | uint32 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.9 | 16-bit value of [802.1ah] B-Tag TPID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.10 | 1 bit B-PCP | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.11 | 1 bit B-DEI | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.6.12 | 12-bit value of [802.1ah] B-VID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.2.8 | 16-bit value of [802.1ad] S-TPID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.3 | eSafe DHCP snooping | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.4 | CM Interface Mask subtype | hexstring | L2VPN I17 | - |
| 43.5.5 | Attachment Group ID (AGI) | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.6 | source attachment individual id (SAII) | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.7 | target attachment individual id (TAII) | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.8 | Upstream User Priority subtype | uint8 | DPoE-MULPIv1.0 C01 | - |
| 43.5.9 | Downstream User Priority Range | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.10 | L2VPN SA-Descriptor Subtypes | hexstring | DPoE-MULPIv1.0 C01 | - |
| 43.5.43 | Vendor Specific L2VPN Subtype | hexstring | L2VPN I17 | - |
| 43.5.12 | Pseudowire Type | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.13 | L2VPN Mode | uint8 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14 | Tag Protocol Identifier (TPID) Translation | encapsulated | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.1 | Upstream TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.2 | Downstream TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.3 | Upstream S-TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.4 | Downstream S-TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.5 | Upstream B-TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.6 | Downstream B-TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.7 | Upstream I-TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.14.8 | Downstream I-TPID Translation | uint16 | DPoE-MULPIv1.0 C01 | - |
| 43.5.15 | L2CP Processing | encapsulated | DPoE-MULPIv1.0 C01 / DPoE-MULPIv2.0 I14 | - |
| 43.5.15.1 | L2CP Tunnel Mode | uint8 | DPoE-MULPIv1.0 C01 / DPoE-MULPIv2.0 I14 | - |
| 43.5.15.2 | L2CP D-MAC Address | macAddress | DPoE-MULPIv1.0 C01 | - |
| 43.5.15.3 | L2CP L2PT D-MAC Address | macAddress | DPoE-MULPIv1.0 C01 / DPoE-MULPIv2.0 I14 | - |
| 43.5.15.4 | L2CP Filter | hexstring | DPoE-MULPIv2.0 I14 | - |
| 43.5.16 | Reserved (formerly DAC) | - | DPoE-MULPIv2.0 I14 | - |
| 43.5.18 | Pseudowire Class | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.19 | Service Delimiter | encapsulated | DPoE-MULPIv2.0 I14 | - |
| 43.5.19.1 | C-VID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.19.2 | S-VID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.19.3 | I-SID | uint32 | DPoE-MULPIv2.0 I14 | - |
| 43.5.19.4 | B-VID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.20 | VSI Encoding | encapsulated | DPoE-MULPIv2.0 I14 | - |
| 43.5.20.1 | VPLS Class | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.20.2 | E-Tree Role | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.20.3 | E-Tree Root VID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.20.4 | E-Tree Leaf VID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.21 | BGP Attribute | encapsulated | DPoE-MULPIv2.0 I14 | - |
| 43.5.21.1 | BGP VPNID | hexstring | DPoE-MULPIv2.0 I14 | - |
| 43.5.21.2 | Route Distinguisher | hexstring | DPoE-MULPIv2.0 I14 | - |
| 43.5.21.3 | Route Target (import) | hexstring | DPoE-MULPIv2.0 I14 | - |
| 43.5.21.4 | Route Target (export) | hexstring | DPoE-MULPIv2.0 I14 | - |
| 43.5.21.5 | CE-ID or VE-ID | uint16 | DPoE-MULPIv2.0 I14 | - |
| 43.5.22 | VPN-SG Attribute | hexstring | DPoE-MULPIv2.0 I14 | - |
| 43.5.23 | Pseudowire Signaling | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.24 | L2VPN SOAM Subtype | encapsulated | L2VPN I17 | - |
| 43.5.24.1 | MEP Configuration | encapsulated | L2VPN I17 | - |
| 43.5.24.1.1 | MD Level | uint8 | L2VPN I17 | - |
| 43.5.24.1.2 | MD Name | string | L2VPN I17 | - |
| 43.5.24.1.3 | MA Name | string | L2VPN I17 | - |
| 43.5.24.1.4 | MEP ID | uint16 | L2VPN I17 | - |
| 43.5.24.2 | Remote MEP Configuration | encapsulated | L2VPN I17 | - |
| 43.5.24.2.1 | MD Level | uint8 | L2VPN I17 | - |
| 43.5.24.2.2 | MD Name | string | L2VPN I17 | - |
| 43.5.24.2.3 | MA Name | string | L2VPN I17 | - |
| 43.5.24.2.4 | MEP ID | uint16 | L2VPN I17 | - |
| 43.5.24.3 | Fault Management Configuration | encapsulated | L2VPN I17 | - |
| 43.5.24.3.1 | Continuity Check Messages | uint8 | L2VPN I17 | - |
| 43.5.24.3.2 | Enable Loopback Reply Messages | uint8 | L2VPN I17 | - |
| 43.5.24.3.3 | Enable Linktrace Messages | uint8 | L2VPN I17 | - |
| 43.5.24.4 | Performance Management Configuration | encapsulated | L2VPN I17 | - |
| 43.5.24.4.1 | Frame Delay Measurement | uint8 | L2VPN I17 | - |
| 43.5.24.4.2 | Frame Loss Measurement | uint8 | L2VPN I17 | - |
| 43.5.25 | Network Timing Profile Reference | uint8 | L2VPN I17 | - |
| 43.5.26 | L2VPN DSID | uint32 | L2VPN I17 | - |
| 43.5.27 | Multipoint Enable/Disable | uint8 | DPoE-MULPIv2.0 I14 | - |
| 43.5.254 | L2VPN Error Encoding | encapsulated | L2VPN I17 | - |
| 43.5.254.1 | L2VPN Errored Parameter | uint8 | L2VPN I17 | - |
| 43.5.254.2 | L2VPN Confirmation Code | uint8 | L2VPN I17 | - |
| 43.5.254.3 | L2VPN Error Message Subtype | string | L2VPN I17 | - |

### 3.4 TLV 45 -- L2VPN DUT Filtering Sub-TLVs

CANN Section 11.1.2.2. Sub-TLVs within TLV 45 (Downstream Unencrypted Traffic Filtering).

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 45.1 | Downstream Unencrypted Traffic (DUT) Control | uint8 | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 45.2 | Downstream Unencrypted Traffic (DUT) CMIM | hexstring | DOCSIS 2.0 · MULPIv4.0 I11 | - |

### 3.5 TLV 65 -- L2VPN MAC Aging Sub-TLVs

CANN Section 11.1.2.3. Sub-TLVs within TLV 65 (L2VPN MAC Aging Encoding).

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 65.1 | L2VPN MAC Aging Mode | uint8 | L2VPN I17 | - |

### 3.6 TLV 24/25/70/71 -- Service Flow Sub-TLVs

CANN Section 11.1.3. These sub-TLVs are shared across Upstream Service Flow (24), Downstream Service Flow (25), Upstream Aggregate Service Flow (70), and Downstream Aggregate Service Flow (71). The applicability column indicates which parent TLVs use each sub-TLV.

#### Common Service Flow Sub-TLVs (24/25/70/71)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.1 | Service Flow Reference or ASF Reference | uint16 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.2 | Service Flow Identifier or ASF Identifier | uint32 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.3 | Service Identifier | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.4 | Service Class Name | string | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.5 | Service Flow Error Encoding | encapsulated | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.5.1 | Errored Parameter | uint8 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.5.2 | Error Code | uint8 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.5.3 | Error Message | string | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.6 | Quality of Service Parameter Set Type | uint8 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.7 | Traffic Priority | uint8 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.8 | Maximum Sustained Traffic Rate | uint32 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9 | Maximum Traffic Burst | uint32 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.10 | Minimum Reserved Traffic Rate | uint32 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.11 | Assumed Minimum Reserved Rate Packet Size | uint16 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12 | Timeout for Active QoS Parameters | uint16 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.13 | Timeout for Admitted QoS Parameters | uint16 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.23 | IP Type Of Service (DSCP) Overwrite | hexstring | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.27 | Peak Traffic Rate | uint32 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.31 | Service Flow Required Attribute Mask | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.32 | Service Flow Forbidden Attribute Mask | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.33 | Service Flow Attribute Aggregation Rule Mask | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.34 | Application Identifier | string | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.35 | Buffer Control | encapsulated | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.35.1 | Minimum Buffer | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.35.2 | Target Buffer | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.35.3 | Maximum Buffer | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.36 | Aggregate Service Flow Reference | uint16 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.37 | Metro Ethernet Service Profile (MESP) Reference | uint16 | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.38 | Serving Group Name | string | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.40 | AQM Encodings | encapsulated | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.40.1 | AQM Disable | uint8 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.40.2 | AQM Latency Target | uint32 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.40.3 | AQM Algorithm | uint8 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.40.4 | Immediate AQM Min Threshold | uint32 | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| x.40.5 | Immediate AQM Range Exponent of Ramp Function | uint8 | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| x.40.6 | Latency Histogram Encodings | hexstring | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| x.41 | Data Rate Unit Setting | uint8 | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.43 | Vendor Specific QoS Parameters | hexstring | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.44 | Guaranteed Grant Interval (GGI) / Service Flow Collection | uint32 | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 / DPoE-MULPIv1.0 C01 | - |

#### ASF-Specific Sub-TLVs (70/71 only)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| 70/71.38.1 | Service Flow to ASF Matching by Application Id | string | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.38.2 | Service Flow to ASF Matching by Service Class Name | string | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.38.3 | Service Flow to ASF Matching by Traffic Priority Range | hexstring | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.1 | Low Latency Service Flow Reference | uint16 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.2 | Low Latency Service Flow Identifier | uint32 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.3 | Classic SF SCN | string | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.4 | Low Latency SF SCN | string | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.5 | AQM Coupling Factor Exponent | uint8 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.6 | Scheduling Weight | uint8 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.7 | Queue Protection Enable | uint8 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.8 | QPLatencyThreshold (CRITICALqL_us) | uint32 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.9 | QPQueuingScoreThreshold (CRITICALqLSCORE_us) | uint32 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |
| 70/71.42.10 | QPDrainRateExponent(LG_AGING) | uint8 | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 | - |

#### TLV 24 Upstream-Only Sub-TLVs

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 24.14 | Maximum Concatenated Burst | uint16 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.15 | Service Flow Scheduling Type | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.16 | Request/Transmission Policy | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.17 | Nominal Polling Interval | uint32 | DPoE-MULPIv1.0 C01 | - |
| 24.18 | Tolerated Poll Jitter | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.19 | Unsolicited Grant Size | uint16 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.20 | Nominal Grant Interval | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.21 | Tolerated Grant Jitter | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.22 | Grants per Interval | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.24 | Unsolicited Grant Time Reference | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.25 | Multiplier to Contention Request Backoff Window | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 24.26 | Multiplier to Number of Bytes Requested | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### TLV 25 Downstream-Only Sub-TLVs

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 25.14 | Maximum Downstream Latency | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 25.15 | Reserved | - | - | - |
| 25.17 | Downstream Resequencing | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

### 3.7 TLV 22/23/60 -- Classification Sub-TLVs

CANN Section 11.1.4. These sub-TLVs are shared across Upstream Packet Classification (22), Downstream Packet Classification (23), and Upstream Drop Packet Classification (60).

#### Common Classification Sub-TLVs

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.1 | Classifier Reference | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.2 | Classifier Identifier | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.3 | Service Flow Reference | uint16 | 22, 23 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.4 | Service Flow Identifier | uint32 | 22, 23 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.5 | Rule Priority | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.6 | Classifier Activation State | uint8 | 22, 23 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.7 | Dynamic Service Change Action | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.8 | Classifier Error Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.8.1 | Errored Parameter | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.8.2 | Error Code | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.8.3 | Error Message | string | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### IPv4/TCP/UDP Classification (x.9)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.9 | IPv4 Packet Classification Encodings / TCP/UDP Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.1 | IPv4 Type of Service Range and Mask | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.2 | IP Protocol | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.3 | IPv4 Source Address | ipv4Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.4 | IPv4 Source Mask | ipv4Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.5 | IPv4 Destination Address | ipv4Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.6 | IPv4 Destination Mask | ipv4Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.7 | TCP/UDP Source Port Start | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.8 | TCP/UDP Source Port End | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.9 | TCP/UDP Destination Port Start | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.9.10 | TCP/UDP Destination Port End | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### Ethernet LLC Classification (x.10)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.10 | Ethernet LLC Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.10.1 | Destination MAC Address | macAddress | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.10.2 | Source MAC Address | macAddress | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.10.3 | Ethertype/DSAP/Mac Type | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.10.4 | Slow Protocol Subtype | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### IEEE 802.1P/Q Classification (x.11)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.11 | IEEE 802.1P/Q Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.11.1 | IEEE 802.1P User Priority | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.11.2 | IEEE 802.1Q VLAN_ID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### IPv6 Classification (x.12)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.12 | IPv6 Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.1 | IPv6 Traffic Class | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.2 | IPv6 Flow Label | uint32 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.3 | IPv6 Next Header Type | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.4 | IPv6 Source Address | ipv6Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.5 | IPv6 Source Prefix Length (bits) | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.6 | IPv6 Destination Address | ipv6Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.12.7 | IPv6 Destination Prefix Length (bits) | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### CM Interface Mask (x.13)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.13 | CM Interface Mask (CMIM) Encoding | hexstring | 22, 23 | L2VPN I17 | - |

#### IEEE 802.1ad S-VLAN Classification (x.14)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.14 | [IEEE 802.1ad] S-VLAN Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.1 | [IEEE 802.1ad] S-TPID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.2 | [IEEE 802.1ad] S-VID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.3 | [IEEE 802.1ad] S-PCP | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.4 | [IEEE 802.1ad] S-DEI | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.5 | [IEEE 802.1ad] C-TPID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.6 | [IEEE 802.1ad] C-VID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.7 | [IEEE 802.1ad] C-PCP | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.8 | [IEEE 802.1ad] C-CFI | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.9 | [IEEE 802.1ad] S-TCI | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.14.10 | [IEEE 802.1ad] C-TCI | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### IEEE 802.1ah I-TAG Classification (x.15)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.15 | [IEEE 802.1ah] I-TAG Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.1 | [IEEE 802.1ah] I-TPID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.2 | [IEEE 802.1ah] I-SID | uint32 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.3 | [IEEE 802.1ah] I-TCI | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.4 | [IEEE 802.1ah] I-PCP | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.5 | [IEEE 802.1ah] I-DEI | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.6 | [IEEE 802.1ah] I-UCA | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.7 | [IEEE 802.1ah] B-TPID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.8 | [IEEE 802.1ah] B-TCI | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.9 | [IEEE 802.1ah] B-PCP | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.10 | [IEEE 802.1ah] B-DEI | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.11 | [IEEE 802.1ah] B-VID | uint16 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.12 | [IEEE 802.1ah] B-DA | macAddress | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.15.13 | [IEEE 802.1ah] B-SA | macAddress | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### ICMPv4/ICMPv6 Classification (x.16)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.16 | ICMPv4/ICMPv6 Packet Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.16.1 | ICMPv4/ICMPv6 Type Start | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.16.2 | ICMPv4/ICMPv6 Type End | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### MPLS Classification (x.17)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.17 | MPLS Classification Encodings | encapsulated | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.17.1 | MPLS TC Bits | uint8 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| x.17.2 | MPLS Label | uint32 | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

#### Vendor-Specific and Extension (x.43)

| Type | CANN Name | Data Type | Applies To | Spec | Reference |
|------|-----------|-----------|------------|------|-----------|
| x.43 | Vendor-Specific Classifier Parameters | hexstring | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 23/60.43.5.1 | VPN Identifier | hexstring | 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 23/60.43.8 | General Extension Information | encapsulated | 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 | - |

### 3.8 TLV 26 -- Payload Header Suppression Sub-TLVs

CANN Section 11.1.5. Sub-TLVs within TLV 26 (Payload Header Suppression).

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 26.1 | Classifier Reference | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.2 | Classifier Identifier | uint16 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.3 | Service Flow Reference | uint16 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.4 | Service Flow Identifier | uint32 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.5 | Dynamic Service Change Action | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.6 | Payload Header Suppression Error Encodings | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.6.1 | Errored Parameter | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.6.2 | Error Code | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.6.3 | Error Message | string | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.7 | Payload Header Suppression Field (PHSF) | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.8 | Payload Header Suppression Index (PHSI) | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.9 | Payload Header Suppression Mask (PHSM) | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.10 | Payload Header Suppression Size (PHSS) | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.11 | Payload Header Suppression Verification (PHSV) | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.13 | Dynamic Bonding Change Action | uint8 | DOCSIS 3.0 · MULPIv4.0 I11 | - |
| 26.43 | Vendor Specific PHS Parameters | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | - |

### 3.9 TLV 53/54 -- SNMP Sub-TLVs

CANN Section 11.1.6. Sub-TLVs within TLV 53 (SNMPv1v2c Coexistence) and TLV 54 (SNMPv3 Access View Configuration).

#### TLV 53 -- SNMPv1v2c Coexistence Sub-TLVs

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 53.1 | SNMPv1v2c Community Name | string | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 53.2 | SNMPv1v2c Transport Address Access | encapsulated | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 53.2.1 | SNMPv1v2c Transport Address | hexstring | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 53.2.2 | SNMPv1v2c Transport Address Mask | hexstring | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 53.3 | SNMPv1v2c Access View Type | uint8 | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 53.4 | SNMPv1v2c Access View Name | string | DOCSIS 2.0 · MULPIv4.0 I11 | - |

#### TLV 54 -- SNMPv3 Access View Configuration Sub-TLVs

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 54.1 | SNMPv3 Access View Name | string | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 54.2 | SNMPv3 Access View Subtree | oid | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 54.3 | SNMPv3 Access View Mask | hexstring | DOCSIS 2.0 · MULPIv4.0 I11 | - |
| 54.4 | SNMPv3 Access View Type | uint8 | DOCSIS 2.0 · MULPIv4.0 I11 | - |

### 3.10 TLV 72 -- MESP Sub-TLVs

CANN Section 11.1.7. Sub-TLVs within TLV 72 (Metro Ethernet Service Profile).

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 72.1 | MESP Reference | uint16 | DPoE-MULPIv2.0 I14 | - |
| 72.2 | MESP Bandwidth Profile (MESP-BP) | encapsulated | DPoE-MULPIv2.0 I14 | - |
| 72.2.1 | MESP-BP Committed Information Rate | uint32 | DPoE-MULPIv2.0 I14 | - |
| 72.2.2 | MESP-BP Committed Burst Size | uint32 | DPoE-MULPIv2.0 I14 | - |
| 72.2.3 | MESP-BP Excess Information Rate | uint32 | DPoE-MULPIv2.0 I14 | - |
| 72.2.4 | MESP-BP Excess Burst Size | uint32 | DPoE-MULPIv2.0 I14 | - |
| 72.2.5 | MESP-BP Coupling Flag | uint8 | DPoE-MULPIv2.0 I14 | - |
| 72.2.6 | MESP-BP Color Mode | encapsulated | DPoE-MULPIv2.0 I14 | - |
| 72.2.6.1 | MESP-BP-CM Color Identification Field | uint8 | DPoE-MULPIv2.0 I14 | - |
| 72.2.6.2 | MESP-BP-CM Color Identification Field Value | hexstring | DPoE-MULPIv2.0 I14 | - |
| 72.2.7 | MESP-BP Color Marking | encapsulated | DPoE-MULPIv2.0 I14 | - |
| 72.2.7.1 | MESP-BP-CR Color Marking Field | uint8 | DPoE-MULPIv2.0 I14 | - |
| 72.2.7.2 | MESP-BP-CR Color Marking Field Value | hexstring | DPoE-MULPIv2.0 I14 | - |
| 72.3 | MESP Name | string | DPoE-MULPIv2.0 I14 | - |

### 3.11 TLV 83 -- L2CP Sub-TLVs

CANN Section 11.1.8. Sub-TLVs within TLV 83 (L2CP Management).

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 83.1 | CMIM | hexstring | DPoE-MULPIv2.0 I14 | - |
| 83.2 | L2CP Mode | uint8 | DPoE-MULPIv2.0 I14 | - |
| 83.3 | L2CP L2PT D-MAC Address | macAddress | DPoE-MULPIv2.0 I14 | - |
| 83.4 | L2CP Filter | hexstring | DPoE-MULPIv2.0 I14 | - |

### 3.12 TLV 202 -- eRouter Sub-TLVs

CANN Section 11.1.9. Sub-TLVs within TLV 202 (eRouter).

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 202.1 | eRouter Initialization Mode Encoding | uint8 | eRouter I22 | - |
| 202.2 | TR-069 Management Server | encapsulated | eRouter I22 | - |
| 202.2.1 | EnableCWMP | uint8 | eRouter I22 | - |
| 202.2.2 | URL Parameter | string | eRouter I22 | - |
| 202.2.3 | Username Parameter | string | eRouter I22 | - |
| 202.2.4 | Password Parameter | string | eRouter I22 | - |
| 202.2.5 | Connection Request Username | string | eRouter I22 | - |
| 202.2.6 | Connection Request Password | string | eRouter I22 | - |
| 202.2.7 | ACSOverride | uint8 | eRouter I22 | - |
| 202.3 | eRouter Initialization Mode Override | uint8 | eRouter I22 | - |
| 202.10 | Router Advertisement (RA) Transmission Interval | uint16 | eRouter I22 | - |
| 202.11 | SNMP MIB Object | hexstring | eRouter I22 | - |
| 202.12 | IP Multicast Configuration Server | ipv4Address | eRouter I22 | - |
| 202.13 | Link-ID Control | uint8 | eRouter I22 | - |
| 202.42 | Topology Mode Encoding | uint8 | eRouter I22 | - |
| 202.43 | Vendor Specific Information | encapsulated | eRouter I22 | - |
| 202.43.8 | Vendor ID Encoding | hexstring | eRouter I22 | - |
| 202.53 | SNMPv1v2c Coexistence Configuration | encapsulated | eRouter I22 | - |
| 202.53.1 | SNMPv1v2c Community Name | string | eRouter I22 | - |
| 202.53.2 | SNMPv1v2c Community Name | encapsulated | eRouter I22 | - |
| 202.53.2.1 | SNMPv1v2c Transport Address | hexstring | eRouter I22 | - |
| 202.53.2.2 | SNMPv1v2c Transport Address Mask | hexstring | eRouter I22 | - |
| 202.53.2.3 | SNMPv1v2c Access View Type | uint8 | eRouter I22 | - |
| 202.53.2.4 | SNMPv1v2c Access View Name | string | eRouter I22 | - |
| 202.54 | SNMPv3 Access View Configuration | encapsulated | eRouter I22 | - |
| 202.54.1 | SNMPv3 Access View Name | string | eRouter I22 | - |
| 202.54.2 | SNMPv3 Access View Subtree | oid | eRouter I22 | - |
| 202.54.3 | SNMPv3 Access View Mask | hexstring | eRouter I22 | - |
| 202.54.4 | SNMPv3 Access View Type | uint8 | eRouter I22 | - |

### 3.13 TLV 219 -- eTEA Sub-TLVs

CANN Section 11.1.10. Sub-TLVs within TLV 219 (eTEA -- TDM Emulation Adapter). All sub-TLVs are defined in CM-SP-TEI.

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 219.8 | eTEA Symbol Clock | uint8 | TEI I06 | - |
| 219.9 | eTEA IWF cfg-encoding | encapsulated | TEI I06 | - |
| 219.9.1 | eTEA PW Index Setting | uint8 | TEI I06 | - |
| 219.9.2 | eTEA Emulation Type | uint8 | TEI I06 | - |
| 219.9.3 | eTEA PW Peer Address | ipv4Address | TEI I06 | - |
| 219.9.4 | eTEA PW Peer Address IPv6 | ipv6Address | TEI I06 | - |
| 219.9.5 | eTEA PW Destination Port | uint16 | TEI I06 | - |
| 219.9.6 | eTEA PW Peer Destination Port | uint16 | TEI I06 | - |
| 219.9.7 | eTEA PW Name | string | TEI I06 | - |
| 219.9.8 | eTEA PW Description | string | TEI I06 | - |
| 219.9.9 | eTEA PW Admin Status | uint8 | TEI I06 | - |
| 219.9.10 | eTEA Status Change Notification Enable | uint8 | TEI I06 | - |
| 219.9.11 | eTEA PW TDM Type | uint8 | TEI I06 | - |
| 219.9.12 | eTEA PW TDM Configuration Table Index | uint8 | TEI I06 | - |
| 219.9.13 | eTEA PW CESoPSNConfiguration Index | uint8 | TEI I06 | - |
| 219.9.14 | eTEA PW RTP SSRC | uint32 | TEI I06 | - |
| 219.9.15 | eTEA PW Peer SSRC | uint32 | TEI I06 | - |
| 219.9.16 | eTEA PW TDM Circulation Map | encapsulated | TEI I06 | - |
| 219.9.16.1 | eTEA PW TDM Port # | uint8 | TEI I06 | - |
| 219.9.16.2 | eTEA PW TDM Timeslot Map | hexstring | TEI I06 | - |
| 219.10 | eTEA PW TDM Configuration Table | encapsulated | TEI I06 | - |
| 219.10.1 | eTEA PW TDM Configuration Table Index | uint8 | TEI I06 | - |
| 219.10.2 | eTEA PW TDM Payload Size | uint16 | TEI I06 | - |
| 219.10.3 | eTEA PW RTP Header Used | uint8 | TEI I06 | - |
| 219.10.5 | eTEA Jitter Buffer Setting | uint8 | TEI I06 | - |
| 219.10.6 | eTEA PW Payload Suppression | uint8 | TEI I06 | - |
| 219.10.7 | eTEA PW LOPS Exit Criteria | uint8 | TEI I06 | - |
| 219.10.8 | eTEA PW LOPS Entrance Criteria | uint8 | TEI I06 | - |
| 219.10.10 | eTEA PW Packet Replace Policy | uint8 | TEI I06 | - |
| 219.10.11 | eTEA PW Packet Loss Window | uint16 | TEI I06 | - |
| 219.10.12 | eTEA PW Excessive Loss Threshold | uint8 | TEI I06 | - |
| 219.10.15 | eTEA PW Severe Loss Threshold | uint8 | TEI I06 | - |
| 219.10.16 | eTEA PW RTP Timestamp Mode | uint8 | TEI I06 | - |
| 219.10.17 | eTEA PW Default Fill Pattern | hexstring | TEI I06 | - |
| 219.10.18 | eTEA PW L Flag Payload Policy | uint8 | TEI I06 | - |
| 219.10.19 | eTEA PW TOS | uint8 | TEI I06 | - |
| 219.10.20 | eTEA PW RTP Payload Type | uint8 | TEI I06 | - |
| 219.10.21 | eTEA PW RTP Peer Payload Type | uint8 | TEI I06 | - |
| 219.10.22 | eTEA PW RTP Timestamp Reference | uint32 | TEI I06 | - |
| 219.10.23 | eTEA PW RTP Peer Timestamp Reference | uint32 | TEI I06 | - |
| 219.10.24 | eTEA PW SRTP Enable | uint8 | TEI I06 | - |
| 219.11 | eTEA SNMP MIB Object | hexstring | TEI I06 | - |
| 219.12 | eTEA SNMP Write-Access Control | hexstring | TEI I06 | - |
| 219.13 | CESoPSN Configuration Table | encapsulated | TEI I06 | - |
| 219.13.1 | PW CESoPSN Config Table Index | uint8 | TEI I06 | - |
| 219.13.2 | Default Idle Pattern | hexstring | TEI I06 | - |
| 219.13.3 | LFlag Policy | uint8 | TEI I06 | - |
| 219.13.4 | RFlag Policy | uint8 | TEI I06 | - |
| 219.13.5 | Remote Defect Policy | uint8 | TEI I06 | - |
| 219.13.6 | LOPS Policy | uint8 | TEI I06 | - |
| 219.13.7 | App Sig TOS | uint8 | TEI I06 | - |
| 219.13.8 | RTP CAS PT | uint8 | TEI I06 | - |
| 219.13.9 | RTP CAS Peer PT | uint8 | TEI I06 | - |
| 219.13.10 | App Sig Idle | hexstring | TEI I06 | - |
| 219.13.11 | App Sig Interval | uint16 | TEI I06 | - |
| 219.13.12 | App Sig Max Interval | uint16 | TEI I06 | - |
| 219.14 | dsx1 Configuration Table | encapsulated | TEI I06 | - |
| 219.14.1 | dsx1 Port ID | uint8 | TEI I06 | - |
| 219.14.2 | dsx1 Line Type | uint8 | TEI I06 | - |
| 219.14.3 | dsx1 Line Coding | uint8 | TEI I06 | - |
| 219.14.4 | dsx1 Circuit ID | string | TEI I06 | - |
| 219.14.5 | dsx1 Loopback Configuration | uint8 | TEI I06 | - |
| 219.14.6 | dsx1 Signal Mode | uint8 | TEI I06 | - |
| 219.14.7 | dsx1 Transmit Clock Source | uint8 | TEI I06 | - |
| 219.14.8 | dsx1 Fdl | uint8 | TEI I06 | - |
| 219.14.9 | dsx1 Line Length | uint16 | TEI I06 | - |
| 219.14.10 | dsx1 Line Status Trap Enable | uint8 | TEI I06 | - |
| 219.14.11 | dsx1 Channelization | uint8 | TEI I06 | - |
| 219.14.12 | dsx1 Line Mode | uint8 | TEI I06 | - |
| 219.14.13 | dsx1 Line Build Out | uint8 | TEI I06 | - |
| 219.43 | eTEA Vendor Specific Extensions | hexstring | TEI I06 | - |
| 219.255 | eTEA End of Text | uint8 | TEI I06 | - |

---

## 4. TLV 43 -- Vendor Specific Information

TLV 43 (DOCSIS Extension Field) is a container TLV that carries vendor-specific or general extension sub-TLVs. Its structure is defined in MULPIv4.0 Annex C Section C.1.1.18.

### Structure

TLV 43 **must** begin with sub-TLV 8 (Vendor ID Encoding), which contains a 3-byte IEEE Organizationally Unique Identifier (OUI). The OUI determines how subsequent sub-TLVs within this TLV 43 instance are interpreted:

- **OUI = 0xFFFFFF** -- General Extension Information. The sub-TLVs that follow are CANN-registered types (43.1 through 43.12), with meanings defined by CableLabs specifications. See Section 3.2 above for the full list.
- **OUI = any other value** -- Vendor Specific. The sub-TLV types and meanings are defined by the vendor identified by the OUI. The CMTS/CM must ignore vendor-specific sub-TLVs from unrecognized vendors.

### Encoding Format

```
Type   = 43 (1 byte)
Length = n  (1 byte, total length of value)
Value  = [ Sub-TLV 8: Vendor ID (T=8, L=3, V=OUI) ]
         [ Sub-TLV ... ]
         [ Sub-TLV ... ]
```

**Example -- General Extension with Load Balancing Policy ID:**

```
Type:   43
Length: 10
  Sub-TLV 8 (Vendor ID):
    Type:   8
    Length: 3
    Value:  0xFF 0xFF 0xFF   (General Extension)
  Sub-TLV 1 (CM Load Balancing Policy ID):
    Type:   1
    Length: 4
    Value:  0x00 0x00 0x00 0x01
```

**Example -- Vendor Specific:**

```
Type:   43
Length: 8
  Sub-TLV 8 (Vendor ID):
    Type:   8
    Length: 3
    Value:  0x00 0x10 0x18   (Broadcom OUI)
  Sub-TLV 1 (vendor-defined):
    Type:   1
    Length: 2
    Value:  <vendor-defined>
```

### General Extension Sub-TLVs (OUI = 0xFFFFFF)

The following sub-TLV types are defined when the Vendor ID is 0xFFFFFF. These are also listed in Section 3.2 with their full sub-TLV hierarchies.

| Type | CANN Name | Data Type | Spec | Reference |
|------|-----------|-----------|------|-----------|
| 43.1 | CM Load Balancing Policy ID | uint32 | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.18.1.1 |
| 43.2 | CM Load Balancing Priority | uint8 | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.18.1.2 |
| 43.3 | CM Load Balancing Group ID | uint32 | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.18.1.3 |
| 43.4 | CM Ranging Class ID Extension | uint16 | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.18.1.4 |
| 43.5 | L2VPN Encoding | encapsulated | L2VPN I17 | C.1.1.18.1.5 |
| 43.6 | Extended CMTS MIC Configuration Setting | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.6 |
| 43.7 | Source Address Verification (SAV) | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.7 |
| 43.8 | Cable Modem Attribute Masks | hexstring | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.8 |
| 43.9 | IP Multicast Join Authorization | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.9 |
| 43.10 | Service Type Identifier | string | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.10 |
| 43.12 | DEMARC Auto-Configuration (DAC) | encapsulated | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.18.1.11 |

### Reference

- MULPIv4.0 C.1.1.18 -- DOCSIS Extension Field
- MULPIv4.0 C.1.1.18.1 -- General Extension Information Encodings
- MULPIv4.0 C.1.1.18.2 -- Vendor-Specific Encodings
- CANN Section 11.1.2 -- TLV 43 Sub-TLV Definitions
