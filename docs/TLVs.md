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

| Type | CANN Name | Length | Cfg File | Spec | Reference |
|------|-----------|--------|----------|------|-----------|
| 0 | Pad | - | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.2 |
| 1 | Downstream Frequency | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.1 |
| 2 | Upstream Channel ID | 1 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.2 |
| 3 | Network Access Control Object | 1 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.3 |
| 4 | DOCSIS 1.0 Class of Service | - | - | DOCSIS 1.0 · MULPIv4.0 I11 | *(deprecated)* |
| 5 | Modem Capabilities Encoding | n | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.1 |
| 6 | CM Message Integrity Check (MIC) | 16 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.5 |
| 7 | CMTS Message Integrity Check (MIC) | 16 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.6 |
| 8 | Vendor ID Encoding | 3 | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.2 |
| 9 | SW Upgrade Filename | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.3 |
| 10 | SNMP Write Access Control | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.4 |
| 11 | SNMP MIB Object | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.5 |
| 12 | Modem IP Address | 4 | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.3 |
| 13 | Service(s) Not Available Response | 3 | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.4 |
| 14 | CPE Ethernet MAC Address | 6 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.6 |
| 15 | Telephone Settings Option | - | - | DOCSIS 1.0 · MULPIv4.0 I11 | *(deprecated)* |
| 16 | *(unassigned)* | - | - | - | - |
| 17 | Baseline Privacy (Security) | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.3.1.1 |
| 18 | Max Number of CPEs | 1 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.7 |
| 19 | TFTP Server Timestamp | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.8 |
| 20 | TFTP Server Provisioned Modem IPv4 Address | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.9 |
| 21 | SW Upgrade IPv4 TFTP Server | 4 | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.7 |
| 22 | Upstream Packet Classification | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.1.1 |
| 23 | Downstream Packet Classification | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.1.3 |
| 24 | Upstream Service Flow | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.2.1 |
| 25 | Downstream Service Flow | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.2.2 |
| 26 | Payload Header Suppression | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.2.3 |
| 27 | HMAC-Digest | 20 | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.4.1 |
| 28 | Maximum Number of Classifiers | 2 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.16 |
| 29 | Privacy Enable | 1 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.17 |
| 30 | Authorization Block | n | No | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.4.2 |
| 31 | Key Sequence Number | 1 | No | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.4.3 |
| 32 | Manufacturer Code Verification Certificate | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.10 |
| 33 | Co-Signer Code Verification Certificate | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.11 |
| 34 | SNMPv3 Kickstart Value | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.9 |
| 35 | Subscriber Mgmt Control | 3 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.19.1 |
| 36 | Subscriber Mgmt CPE IPv4 List | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.19.2 |
| 37 | Subscriber Mgmt Filter Groups | 8 | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.1.19.4 |
| 38 | SNMPv3 Notification Receiver | n | Yes | DOCSIS 1.1 · MULPIv4.0 I11 | C.1.2.12 |
| 39 | Enable 2.0 Mode | 1 | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.20 |
| 40 | Enable Test Modes | 1 | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.20 |
| 41 | Downstream Channel List | n | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.22 |
| 42 | Static Multicast MAC Address | 6 | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.23 |
| 43 | DOCSIS Extension Field | n | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.1.18 |
| 44 | Vendor Specific Capabilities | n | No | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.3.5 |
| 45 | Downstream Unencrypted Traffic (DUT) Filtering | n | Yes | DOCSIS 2.0 · MULPIv4.0 I11 | C.1.1.24 |
| 46 | Transmit Channel Configuration (TCC) | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.1 |
| 47 | Service Flow SID Cluster Assignment | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.2 |
| 48 | Receive Channel Profile | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.3.1 |
| 49 | Receive Channel Configuration | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.3.1 |
| 50 | DSID Encodings | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.3.9 |
| 51 | Security Association Encoding | n | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.5 |
| 52 | Initializing Channel Timeout | 2 | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.5.6 |
| 53 | SNMPv1v2c Coexistence | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.13 |
| 54 | SNMPv3 Access View Configuration | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.14 |
| 55 | SNMP CPE Access Control | 1 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.15 |
| 56 | Channel Assignment Configuration Settings | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.25 |
| 57 | CM Initialization Reason | 1 | No | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.3.6 |
| 58 | SW Upgrade IPv6 TFTP Server | 16 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.8 |
| 59 | TFTP Server Provisioned Modem IPv6 Address | 16 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.10 |
| 60 | Upstream Drop Packet Classification | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.2.1.2 |
| 61 | Subscriber Mgmt CPE IPv6 Prefix List | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.19.3 |
| 62 | Upstream Drop Classifier Group ID | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.26 |
| 63 | Subscriber Mgmt Control Max CPE IPv6 Prefix | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.19.5 |
| 64 | CMTS Static Multicast Session Encoding | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.27 |
| 65 | L2VPN MAC Aging Encoding | n | Yes | L2VPN I17 | C.1.1.28 |
| 66 | Management Event Control Encoding | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.16 |
| 67 | Subscriber Mgmt CPE IPv6 List | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.19.6 |
| 68 | Default Upstream Target Buffer Configuration | 2 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.17 |
| 69 | MAC Address Learning Control Encoding | 1 | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.2.18 |
| 70 | Upstream Aggregate Service Flow Encodings | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.3 |
| 71 | Downstream Aggregate Service Flow Encodings | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.4 |
| 72 | Metro Ethernet Service Profile | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 / DPoE-MULPIv2.0 I14 | C.2.2.12 |
| 73 | Network Timing Profile | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 / DPoE-MULPIv2.0 I14 | C.1.2.19 |
| 74 | Energy Management Parameter Encoding | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.1.1.30 |
| 75 | Energy Mgt. Mode Indicator | 1 | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.4.4 |
| 76 | CM Upstream AQM Disable | 1 | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.20 |
| 77 | DOCSIS Time Protocol Encoding | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.6 |
| 78 | Energy Management Identifier List for CM | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.1.30.4 |
| 79 | UNI Control Encoding | n | Yes | DOCSIS 3.0 · MULPIv4.0 I11 | C.3.3 |
| 80 | Energy Management -- DOCSIS Light Sleep Encodings | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.4.5 |
| 81 | Manufacturer CVC Chain | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.21 |
| 82 | Co-signer CVC Chain | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.22 |
| 83 | L2CP Management | 1 | Yes | DPoE-MULPIv2.0 I14 | C.1.1.31 |
| 84 | Diplexer Band Edge | 9 | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.2.23 |
| 85 | FDX Transmission Group Assignment | n | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.4.6 |
| 86 | FDX Reset | 1 | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.4.7 |
| 87 | CM Echo Cancellation Training Control | n | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.7 |
| 88 | QoS Framework for DOCSIS Encodings | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.8 |
| 89 | Extended SID Cluster Assignment | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.3.1.70 |
| 90 | Primary Service Flow Indicator | n | No | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.3.7 |
| 91 | Low Latency Disable | 1 | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.1.32 |
| 92 | Distributed HQoS Enable | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.1.33 |
| 93 | Upstream Enhanced HQoS ASF | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.5 |
| 94 | Downstream Enhanced HQoS ASF | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.2.2.6 |
| 95 | DHQoS ASF SID Bundle Assignment | n | Yes | DOCSIS 3.1 · MULPIv4.0 I11 | C.1.5.8 |
| 96 | Advanced Diplexer Band Edge | n | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.2.24 |
| 97 | Advanced Band Plan Support | 1 | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.2.25 |
| 98 | DOCSIS Sync Capabilities | n | No | SYNC I03 | [SYNC] |
| 99 | DOCSIS CM System Information | n | No | SYNC I03 | [SYNC] |
| 100 | Sync DSID Assignment | n | No | SYNC I03 | [SYNC] |
| 101 | DOCSIS Sync Configurations | n | Yes | SYNC I03 | [SYNC] |
| 102 | PTP Address Configurations | n | Yes | SYNC I03 | [SYNC] |
| 103 | CM SSH Server Configuration Settings | n | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.3.1.2 |
| 104 | Security Configuration Settings | n | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.3.1.3 |
| 105 | Extended Modem Capabilities | n | No | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.3.1 |
| 106 | FDX Downstream Upper Band Edge | 2 | Yes | DOCSIS 4.0 · MULPIv4.0 I11 | C.1.2.26 |

> **Note:** Types 16 and 107--200 are reserved/unassigned.

### eSAFE Types (201--231)

| Type | CANN Name | Length | Cfg File | Spec | Reference |
|------|-----------|--------|----------|------|-----------|
| 201 | ePS | n | Yes | eDOCSIS I31 | [eDOCSIS] |
| 202 | eRouter | n | Yes | eRouter I22 | [eRouter] |
| 203--215 | *(Reserved)* | - | - | - | - |
| 216 | eMTA | n | Yes | PacketCable PROV1.5 C01 | [PacketCable] |
| 217 | eSTB | n | Yes | DSG I25 | [DSG] |
| 218 | *(Reserved)* | - | - | - | - |
| 219 | eTEA | n | Yes | TEI I06 | [TEI] |
| 220 | eDVA | n | Yes | PacketCable RST-E-DVA C01 | [PacketCable 2.0] |
| 221 | eSG | n | Yes | eDOCSIS I31 | [eDOCSIS] |
| 222--231 | *(Reserved)* | - | - | - | - |

> **Note:** Types 232--254 are reserved/unassigned.

### Type 255

| Type | CANN Name | Length | Cfg File | Spec | Reference |
|------|-----------|--------|----------|------|-----------|
| 255 | End-of-Data | - | Yes | DOCSIS 1.0 · MULPIv4.0 I11 | C.1.2.1 |

---

## 3. Sub-TLV Sections

### 3.1 TLV 5 -- Modem Capabilities Encoding Sub-TLVs

CANN Section 11.1.1. These sub-TLVs are carried within TLV 5 (Modem Capabilities Encoding) and are sent by the CM to the CMTS during registration to advertise device capabilities.

| Type | CANN Name | Spec |
|------|-----------|------|
| 5.1 | Concatenation Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.2 | DOCSIS Version | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.3 | Fragmentation Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.4 | Payload Header Suppression Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.5 | IGMP Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.6 | Privacy Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.7 | Downstream SAID Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.8 | Upstream Service Flow Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.9 | Optional Filtering Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.10 | Transmit Pre-Equalizer Taps | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.11 | Number of Transmit Pre-Equalizer Taps | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.12 | DCC Support | DOCSIS 1.1 · MULPIv4.0 I11 |
| 5.13 | IP Filters Support | DOCSIS 2.0 · MULPIv4.0 I11 |
| 5.14 | LLC Filters Support | DOCSIS 2.0 · MULPIv4.0 I11 |
| 5.15 | Expanded Unicast SID Space | DOCSIS 2.0 · MULPIv4.0 I11 |
| 5.16 | Ranging Hold-off Support | DOCSIS 2.0 · MULPIv4.0 I11 |
| 5.17 | L2VPN Capability | L2VPN I17 |
| 5.18 | L2VPN eSAFE Host Capability | L2VPN I17 |
| 5.19 | Downstream Unencrypted Traffic (DUT) Filtering | L2VPN I17 |
| 5.20 | Upstream Frequency Range Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.21 | Upstream SC-QAM Symbol Rate Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.22 | Selectable Active Code Mode 2 Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.23 | Code Hopping Mode 2 Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.24 | Multiple Transmit SC-QAM Channel Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.25 | 5.12 Msps UpstreamTransmit SC-QAM Channel Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.26 | 2.56 Msps Upstream Transmit SC-QAM Channel Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.27 | Total SID Cluster Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.28 | SID Clusters per Service Flow Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.29 | Multiple Receive SC-QAM Channel Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.30 | Total Downstream Service ID (DSID) Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.31 | Resequencing Downstream Service ID (DSID) Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.32 | Multicast Downstream Service ID (DSID) Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.33 | Multicast DSID Forwarding | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.34 | Frame Control Type Forwarding Capability | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.35 | DPV Capability | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.36 | Unsolicited Grant Service/Upstream Service Flow Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.37 | MAP and UCD Receipt Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.38 | Upstream Drop Classifier Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.39 | IPv6 Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.40 | Extended Upstream Transmit Power Capability | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.41 | Optional 802.1ad, 802.1ah, MPLS Classification Support | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.42 | D-ONU Capabilities | DPoE-MULPIv1.0 C01 |
| 5.42.1 | DPoE/G Version Number | DPoE-MULPIv1.0 C01 |
| 5.42.2 | Number of Unicast LLIDs | DPoE-MULPIv1.0 C01 |
| 5.42.3 | Number of Multicast LLIDs | DPoE-MULPIv2.0 I14 |
| 5.42.4 | MESP Support | DPoE-MULPIv2.0 I14 |
| 5.42.5 | Number of D-ONU Ports | DPoE-MULPIv2.0 I14 |
| 5.42.6 | PON Data Rate Support | DPoE-MULPIv2.0 I14 |
| 5.42.7 | Service OAM | DPoE-MULPIv2.0 I14 |
| 5.42.10 | Number of T-CONTs Supported | DPoG 1.0 |
| 5.42.11 | Total Number of (X)GEM Ports Supported | DPoG 1.0 |
| 5.43 | Reserved | - |
| 5.44 | Energy Management Capabilities | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.45 | C-DOCSIS Capability Encoding | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.46 | CM-STATUS-ACK | DOCSIS 3.0 · MULPIv4.0 I11 |
| 5.47 | Energy Management Preferences | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.48 | Extended Packet Length Support Capability | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.49 | Multiple Receive OFDM Channel Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.50 | Multiple Transmit OFDMA Channel Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.51 | Downstream OFDM Profile Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.52 | Downstream OFDM Channel Subcarrier QAM Modulation Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.53 | Upstream OFDMA Channel Subcarrier QAM Modulation Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.54 | Downstream Lower Band Edge Configuration | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.55 | Downstream Upper Band Edge Configuration | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.56 | Diplexer Upstream Upper Band Edge Configuration | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.57 | DOCSIS Time Protocol Mode | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.58 | DOCSIS Time Protocol Performance Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.59 | Pmax | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.60 | Diplexer Downstream Lower Band Edge Options | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.61 | Diplexer Downstream Upper Band Edge Options | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.62 | Diplexer Upstream Upper Band Edge Options | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.63 | Advanced Band Plan Capability | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.64 | FDX DS State Lock -- Deprecated | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.65 | FDX Switching Software Timing Uncertainty | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.66 | FDX DS to US Switching Time | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.67 | FDX US to DS Switching Time | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.68 | -- | - |
| 5.69 | CWT RxMER Measurement Convergence Time | - |
| 5.70 | -- | - |
| 5.71 | -- | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.72 | t-ds-reacquisition capability | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.73 | CWT Simultaneous Data Transmission Capability | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.74 | Extended Service Flow SID Cluster Assignments Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.75 | Echo Cancelling RBA Sub-band Direction Sets Supported | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.76 | Low Latency Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.77 | Absolute Queue-Depth Request Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.78 | Distributed HQoS Support | DOCSIS 3.1 · MULPIv4.0 I11 |
| 5.79 | Advanced Downstream Lower Band Edge Configuration | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.80 | Advanced Downstream Upper Band Edge Configuration | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.81 | Advanced Diplexer Upstream Upper Band Edge Configuration | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.82 | Advanced Diplexer Downstream Lower Band Edge Options List | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.83 | Advanced Diplexer Downstream Upper Band Edge Options List | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.84 | Advanced Diplexer Upstream Upper Band Edge Options List | DOCSIS 4.0 · MULPIv4.0 I11 |
| 5.85 | Extended Power Options | DOCSIS 4.0 · MULPIv4.0 I11 |

### 3.2 TLV 43 -- DOCSIS Extension Field Sub-TLVs (General Extension)

CANN Section 11.1.2. When TLV 43 carries OUI = 0xFFFFFF (General Extension Information), the following sub-TLVs are defined. See also [Section 4](#4-tlv-43----vendor-specific-information) for the full TLV 43 structure.

| Type | CANN Name | Reference |
|------|-----------|-----------|
| 43.1 | CM Load Balancing Policy ID | C.1.1.18.1.1 |
| 43.2 | CM Load Balancing Priority | C.1.1.18.1.2 |
| 43.3 | CM Load Balancing Group ID | C.1.1.18.1.3 |
| 43.4 | CM Ranging Class ID Extension | C.1.1.18.1.4 |
| 43.5 | L2VPN Encoding | C.1.1.18.1.5 |
| 43.6 | Extended CMTS MIC Configuration Setting | C.1.1.18.1.6 |
| 43.7 | Source Address Verification (SAV) | C.1.1.18.1.7 |
| 43.8 | Cable Modem Attribute Masks | C.1.1.18.1.8 |
| 43.9 | IP Multicast Join Authorization | C.1.1.18.1.9 |
| 43.10 | Service Type Identifier | C.1.1.18.1.10 |
| 43.12 | DEMARC Auto-Configuration (DAC) | C.1.1.18.1.11 |

> **Note:** Sub-type 43.11 is unassigned.

### 3.3 TLV 43.5 -- L2VPN Encoding Sub-TLVs

CANN Section 11.1.2.1. These sub-TLVs are carried within TLV 43.5 (L2VPN Encoding) inside the General Extension Information field.

| Type | CANN Name | Spec |
|------|-----------|------|
| 43.5.1 | VPN Identifier | L2VPN I17 |
| 43.5.2 | NSI encapsulation format | DPoE-MULPIv1.0 C01 |
| 43.5.2.1 | Other | DPoE-MULPIv1.0 C01 |
| 43.5.2.2 | IEEE 802.1Q | DPoE-MULPIv1.0 C01 |
| 43.5.2.3 | IEEE 802.1ad | DPoE-MULPIv1.0 C01 |
| 43.5.2.4 | MPLS PW | DPoE-MULPIv1.0 C01 |
| 43.5.2.4.1 | MPLS Pseudowire ID | DPoE-MULPIv2.0 I14 |
| 43.5.2.4.2 | MPLS Peer IP address | DPoE-MULPIv2.0 I14 |
| 43.5.2.4.3 | Pseudowire Type | DPoE-MULPIv1.0 C01 |
| 43.5.2.4.4 | MPLS Backup Pseudowire ID | DPoE-MULPIv2.0 I14 |
| 43.5.2.4.5 | MPLS Backup Peer IP address | DPoE-MULPIv2.0 I14 |
| 43.5.2.5 | L2TPv3 Peer | DPoE-MULPIv1.0 C01 |
| 43.5.2.6 | IEEE 802.1ah Encapsulation | DPoE-MULPIv1.0 C01 |
| 43.5.2.6.1 | IEEE 802.1ah Backbone Service Instance Tag (I-Tag) TCI | DPoE-MULPIv1.0 C01 |
| 43.5.2.6.2 | IEEE 802.1ah Destination Backbone Edge Bridge (BEB) MAC Address (B-DA) | DPoE-MULPIv1.0 C01 |
| 43.5.2.6.3 | 16-bit value of [802.1ah] B-Tag TCI | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.4 | 16-bit value of [802.1ah] I-Tag TPID | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.5 | 3 bit I-PCP | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.6 | 1 bit I-DEI | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.7 | 1 bit I-UCA | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.8 | 24-bit value of [802.1ah] I-SID Backbone Service Instance Identifier | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.9 | 16-bit value of [802.1ah] B-Tag TPID | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.10 | 1 bit B-PCP | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.11 | 1 bit B-DEI | DPoE-MULPIv2.0 I14 |
| 43.5.2.6.12 | 12-bit value of [802.1ah] B-VID | DPoE-MULPIv2.0 I14 |
| 43.5.2.8 | 16-bit value of [802.1ad] S-TPID | DPoE-MULPIv2.0 I14 |
| 43.5.3 | eSafe DHCP snooping | DPoE-MULPIv1.0 C01 |
| 43.5.4 | CM Interface Mask subtype | L2VPN I17 |
| 43.5.5 | Attachment Group ID (AGI) | DPoE-MULPIv1.0 C01 |
| 43.5.6 | source attachment individual id (SAII) | DPoE-MULPIv1.0 C01 |
| 43.5.7 | target attachment individual id (TAII) | DPoE-MULPIv1.0 C01 |
| 43.5.8 | Upstream User Priority subtype | DPoE-MULPIv1.0 C01 |
| 43.5.9 | Downstream User Priority Range | DPoE-MULPIv1.0 C01 |
| 43.5.10 | L2VPN SA-Descriptor Subtypes | DPoE-MULPIv1.0 C01 |
| 43.5.43 | Vendor Specific L2VPN Subtype | L2VPN I17 |
| 43.5.12 | Pseudowire Type | DPoE-MULPIv2.0 I14 |
| 43.5.13 | L2VPN Mode | DPoE-MULPIv1.0 C01 |
| 43.5.14 | Tag Protocol Identifier (TPID) Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.1 | Upstream TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.2 | Downstream TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.3 | Upstream S-TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.4 | Downstream S-TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.5 | Upstream B-TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.6 | Downstream B-TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.7 | Upstream I-TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.14.8 | Downstream I-TPID Translation | DPoE-MULPIv1.0 C01 |
| 43.5.15 | L2CP Processing | DPoE-MULPIv1.0 C01 / DPoE-MULPIv2.0 I14 |
| 43.5.15.1 | L2CP Tunnel Mode | DPoE-MULPIv1.0 C01 / DPoE-MULPIv2.0 I14 |
| 43.5.15.2 | L2CP D-MAC Address | DPoE-MULPIv1.0 C01 |
| 43.5.15.3 | L2CP L2PT D-MAC Address | DPoE-MULPIv1.0 C01 / DPoE-MULPIv2.0 I14 |
| 43.5.15.4 | L2CP Filter | DPoE-MULPIv2.0 I14 |
| 43.5.16 | Reserved (formerly DAC) | DPoE-MULPIv2.0 I14 |
| 43.5.18 | Pseudowire Class | DPoE-MULPIv2.0 I14 |
| 43.5.19 | Service Delimiter | DPoE-MULPIv2.0 I14 |
| 43.5.19.1 | C-VID | DPoE-MULPIv2.0 I14 |
| 43.5.19.2 | S-VID | DPoE-MULPIv2.0 I14 |
| 43.5.19.3 | I-SID | DPoE-MULPIv2.0 I14 |
| 43.5.19.4 | B-VID | DPoE-MULPIv2.0 I14 |
| 43.5.20 | VSI Encoding | DPoE-MULPIv2.0 I14 |
| 43.5.20.1 | VPLS Class | DPoE-MULPIv2.0 I14 |
| 43.5.20.2 | E-Tree Role | DPoE-MULPIv2.0 I14 |
| 43.5.20.3 | E-Tree Root VID | DPoE-MULPIv2.0 I14 |
| 43.5.20.4 | E-Tree Leaf VID | DPoE-MULPIv2.0 I14 |
| 43.5.21 | BGP Attribute | DPoE-MULPIv2.0 I14 |
| 43.5.21.1 | BGP VPNID | DPoE-MULPIv2.0 I14 |
| 43.5.21.2 | Route Distinguisher | DPoE-MULPIv2.0 I14 |
| 43.5.21.3 | Route Target (import) | DPoE-MULPIv2.0 I14 |
| 43.5.21.4 | Route Target (export) | DPoE-MULPIv2.0 I14 |
| 43.5.21.5 | CE-ID or VE-ID | DPoE-MULPIv2.0 I14 |
| 43.5.22 | VPN-SG Attribute | DPoE-MULPIv2.0 I14 |
| 43.5.23 | Pseudowire Signaling | DPoE-MULPIv2.0 I14 |
| 43.5.24 | L2VPN SOAM Subtype | L2VPN I17 |
| 43.5.24.1 | MEP Configuration | L2VPN I17 |
| 43.5.24.1.1 | MD Level | L2VPN I17 |
| 43.5.24.1.2 | MD Name | L2VPN I17 |
| 43.5.24.1.3 | MA Name | L2VPN I17 |
| 43.5.24.1.4 | MEP ID | L2VPN I17 |
| 43.5.24.2 | Remote MEP Configuration | L2VPN I17 |
| 43.5.24.2.1 | MD Level | L2VPN I17 |
| 43.5.24.2.2 | MD Name | L2VPN I17 |
| 43.5.24.2.3 | MA Name | L2VPN I17 |
| 43.5.24.2.4 | MEP ID | L2VPN I17 |
| 43.5.24.3 | Fault Management Configuration | L2VPN I17 |
| 43.5.24.3.1 | Continuity Check Messages | L2VPN I17 |
| 43.5.24.3.2 | Enable Loopback Reply Messages | L2VPN I17 |
| 43.5.24.3.3 | Enable Linktrace Messages | L2VPN I17 |
| 43.5.24.4 | Performance Management Configuration | L2VPN I17 |
| 43.5.24.4.1 | Frame Delay Measurement | L2VPN I17 |
| 43.5.24.4.2 | Frame Loss Measurement | L2VPN I17 |
| 43.5.25 | Network Timing Profile Reference | L2VPN I17 |
| 43.5.26 | L2VPN DSID | L2VPN I17 |
| 43.5.27 | Multipoint Enable/Disable | DPoE-MULPIv2.0 I14 |
| 43.5.254 | L2VPN Error Encoding | L2VPN I17 |
| 43.5.254.1 | L2VPN Errored Parameter | L2VPN I17 |
| 43.5.254.2 | L2VPN Confirmation Code | L2VPN I17 |
| 43.5.254.3 | L2VPN Error Message Subtype | L2VPN I17 |

### 3.4 TLV 45 -- L2VPN DUT Filtering Sub-TLVs

CANN Section 11.1.2.2. Sub-TLVs within TLV 45 (Downstream Unencrypted Traffic Filtering).

| Type | CANN Name | Spec |
|------|-----------|------|
| 45.1 | Downstream Unencrypted Traffic (DUT) Control | DOCSIS 2.0 · MULPIv4.0 I11 |
| 45.2 | Downstream Unencrypted Traffic (DUT) CMIM | DOCSIS 2.0 · MULPIv4.0 I11 |

### 3.5 TLV 65 -- L2VPN MAC Aging Sub-TLVs

CANN Section 11.1.2.3. Sub-TLVs within TLV 65 (L2VPN MAC Aging Encoding).

| Type | CANN Name | Spec |
|------|-----------|------|
| 65.1 | L2VPN MAC Aging Mode | L2VPN I17 |

### 3.6 TLV 24/25/70/71 -- Service Flow Sub-TLVs

CANN Section 11.1.3. These sub-TLVs are shared across Upstream Service Flow (24), Downstream Service Flow (25), Upstream Aggregate Service Flow (70), and Downstream Aggregate Service Flow (71). The applicability column indicates which parent TLVs use each sub-TLV.

#### Common Service Flow Sub-TLVs (24/25/70/71)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.1 | Service Flow Reference or ASF Reference | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.2 | Service Flow Identifier or ASF Identifier | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.3 | Service Identifier | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.4 | Service Class Name | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.5 | Service Flow Error Encoding | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.5.1 | Errored Parameter | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.5.2 | Error Code | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.5.3 | Error Message | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.6 | Quality of Service Parameter Set Type | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.7 | Traffic Priority | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.8 | Maximum Sustained Traffic Rate | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9 | Maximum Traffic Burst | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.10 | Minimum Reserved Traffic Rate | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.11 | Assumed Minimum Reserved Rate Packet Size | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12 | Timeout for Active QoS Parameters | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.13 | Timeout for Admitted QoS Parameters | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.23 | IP Type Of Service (DSCP) Overwrite | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.27 | Peak Traffic Rate | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.31 | Service Flow Required Attribute Mask | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.32 | Service Flow Forbidden Attribute Mask | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.33 | Service Flow Attribute Aggregation Rule Mask | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.34 | Application Identifier | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.35 | Buffer Control | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.35.1 | Minimum Buffer | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.35.2 | Target Buffer | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.35.3 | Maximum Buffer | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.36 | Aggregate Service Flow Reference | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.37 | Metro Ethernet Service Profile (MESP) Reference | 24, 25, 70, 71 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.38 | Serving Group Name | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.40 | AQM Encodings | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.40.1 | AQM Disable | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.40.2 | AQM Latency Target | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.40.3 | AQM Algorithm | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.40.4 | Immediate AQM Min Threshold | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 |
| x.40.5 | Immediate AQM Range Exponent of Ramp Function | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 |
| x.40.6 | Latency Histogram Encodings | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 |
| x.41 | Data Rate Unit Setting | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.43 | Vendor Specific QoS Parameters | 24, 25 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.44 | Guaranteed Grant Interval (GGI) / Service Flow Collection | 24, 25 | DOCSIS 3.1 · MULPIv4.0 I11 / DPoE-MULPIv1.0 C01 |

#### ASF-Specific Sub-TLVs (70/71 only)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| 70/71.38.1 | Service Flow to ASF Matching by Application Id | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.38.2 | Service Flow to ASF Matching by Service Class Name | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.38.3 | Service Flow to ASF Matching by Traffic Priority Range | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.1 | Low Latency Service Flow Reference | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.2 | Low Latency Service Flow Identifier | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.3 | Classic SF SCN | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.4 | Low Latency SF SCN | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.5 | AQM Coupling Factor Exponent | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.6 | Scheduling Weight | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.7 | Queue Protection Enable | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.8 | QPLatencyThreshold (CRITICALqL_us) | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.9 | QPQueuingScoreThreshold (CRITICALqLSCORE_us) | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |
| 70/71.42.10 | QPDrainRateExponent(LG_AGING) | 70, 71 | DOCSIS 3.1 · MULPIv4.0 I11 |

#### TLV 24 Upstream-Only Sub-TLVs

| Type | CANN Name | Spec |
|------|-----------|------|
| 24.14 | Maximum Concatenated Burst | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.15 | Service Flow Scheduling Type | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.16 | Request/Transmission Policy | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.17 | Nominal Polling Interval | DPoE-MULPIv1.0 C01 |
| 24.18 | Tolerated Poll Jitter | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.19 | Unsolicited Grant Size | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.20 | Nominal Grant Interval | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.21 | Tolerated Grant Jitter | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.22 | Grants per Interval | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.24 | Unsolicited Grant Time Reference | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.25 | Multiplier to Contention Request Backoff Window | DOCSIS 3.0 · MULPIv4.0 I11 |
| 24.26 | Multiplier to Number of Bytes Requested | DOCSIS 3.0 · MULPIv4.0 I11 |

#### TLV 25 Downstream-Only Sub-TLVs

| Type | CANN Name | Spec |
|------|-----------|------|
| 25.14 | Maximum Downstream Latency | DOCSIS 3.0 · MULPIv4.0 I11 |
| 25.15 | Reserved | - |
| 25.17 | Downstream Resequencing | DOCSIS 3.0 · MULPIv4.0 I11 |

### 3.7 TLV 22/23/60 -- Classification Sub-TLVs

CANN Section 11.1.4. These sub-TLVs are shared across Upstream Packet Classification (22), Downstream Packet Classification (23), and Upstream Drop Packet Classification (60).

#### Common Classification Sub-TLVs

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.1 | Classifier Reference | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.2 | Classifier Identifier | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.3 | Service Flow Reference | 22, 23 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.4 | Service Flow Identifier | 22, 23 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.5 | Rule Priority | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.6 | Classifier Activation State | 22, 23 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.7 | Dynamic Service Change Action | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.8 | Classifier Error Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.8.1 | Errored Parameter | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.8.2 | Error Code | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.8.3 | Error Message | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### IPv4/TCP/UDP Classification (x.9)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.9 | IPv4 Packet Classification Encodings / TCP/UDP Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.1 | IPv4 Type of Service Range and Mask | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.2 | IP Protocol | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.3 | IPv4 Source Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.4 | IPv4 Source Mask | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.5 | IPv4 Destination Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.6 | IPv4 Destination Mask | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.7 | TCP/UDP Source Port Start | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.8 | TCP/UDP Source Port End | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.9 | TCP/UDP Destination Port Start | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.9.10 | TCP/UDP Destination Port End | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### Ethernet LLC Classification (x.10)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.10 | Ethernet LLC Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.10.1 | Destination MAC Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.10.2 | Source MAC Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.10.3 | Ethertype/DSAP/Mac Type | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.10.4 | Slow Protocol Subtype | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### IEEE 802.1P/Q Classification (x.11)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.11 | IEEE 802.1P/Q Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.11.1 | IEEE 802.1P User Priority | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.11.2 | IEEE 802.1Q VLAN_ID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### IPv6 Classification (x.12)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.12 | IPv6 Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.1 | IPv6 Traffic Class | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.2 | IPv6 Flow Label | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.3 | IPv6 Next Header Type | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.4 | IPv6 Source Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.5 | IPv6 Source Prefix Length (bits) | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.6 | IPv6 Destination Address | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.12.7 | IPv6 Destination Prefix Length (bits) | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### CM Interface Mask (x.13)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.13 | CM Interface Mask (CMIM) Encoding | 22, 23 | L2VPN I17 |

#### IEEE 802.1ad S-VLAN Classification (x.14)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.14 | [IEEE 802.1ad] S-VLAN Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.1 | [IEEE 802.1ad] S-TPID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.2 | [IEEE 802.1ad] S-VID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.3 | [IEEE 802.1ad] S-PCP | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.4 | [IEEE 802.1ad] S-DEI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.5 | [IEEE 802.1ad] C-TPID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.6 | [IEEE 802.1ad] C-VID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.7 | [IEEE 802.1ad] C-PCP | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.8 | [IEEE 802.1ad] C-CFI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.9 | [IEEE 802.1ad] S-TCI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.14.10 | [IEEE 802.1ad] C-TCI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### IEEE 802.1ah I-TAG Classification (x.15)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.15 | [IEEE 802.1ah] I-TAG Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.1 | [IEEE 802.1ah] I-TPID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.2 | [IEEE 802.1ah] I-SID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.3 | [IEEE 802.1ah] I-TCI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.4 | [IEEE 802.1ah] I-PCP | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.5 | [IEEE 802.1ah] I-DEI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.6 | [IEEE 802.1ah] I-UCA | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.7 | [IEEE 802.1ah] B-TPID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.8 | [IEEE 802.1ah] B-TCI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.9 | [IEEE 802.1ah] B-PCP | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.10 | [IEEE 802.1ah] B-DEI | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.11 | [IEEE 802.1ah] B-VID | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.12 | [IEEE 802.1ah] B-DA | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.15.13 | [IEEE 802.1ah] B-SA | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### ICMPv4/ICMPv6 Classification (x.16)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.16 | ICMPv4/ICMPv6 Packet Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.16.1 | ICMPv4/ICMPv6 Type Start | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.16.2 | ICMPv4/ICMPv6 Type End | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### MPLS Classification (x.17)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.17 | MPLS Classification Encodings | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.17.1 | MPLS TC Bits | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| x.17.2 | MPLS Label | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

#### Vendor-Specific and Extension (x.43)

| Type | CANN Name | Applies To | Spec |
|------|-----------|------------|------|
| x.43 | Vendor-Specific Classifier Parameters | 22, 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| 23/60.43.5.1 | VPN Identifier | 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |
| 23/60.43.8 | General Extension Information | 23, 60 | DOCSIS 3.0 · MULPIv4.0 I11 |

### 3.8 TLV 26 -- Payload Header Suppression Sub-TLVs

CANN Section 11.1.5. Sub-TLVs within TLV 26 (Payload Header Suppression).

| Type | CANN Name | Spec |
|------|-----------|------|
| 26.1 | Classifier Reference | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.2 | Classifier Identifier | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.3 | Service Flow Reference | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.4 | Service Flow Identifier | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.5 | Dynamic Service Change Action | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.6 | Payload Header Suppression Error Encodings | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.6.1 | Errored Parameter | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.6.2 | Error Code | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.6.3 | Error Message | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.7 | Payload Header Suppression Field (PHSF) | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.8 | Payload Header Suppression Index (PHSI) | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.9 | Payload Header Suppression Mask (PHSM) | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.10 | Payload Header Suppression Size (PHSS) | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.11 | Payload Header Suppression Verification (PHSV) | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.13 | Dynamic Bonding Change Action | DOCSIS 3.0 · MULPIv4.0 I11 |
| 26.43 | Vendor Specific PHS Parameters | DOCSIS 3.0 · MULPIv4.0 I11 |

### 3.9 TLV 53/54 -- SNMP Sub-TLVs

CANN Section 11.1.6. Sub-TLVs within TLV 53 (SNMPv1v2c Coexistence) and TLV 54 (SNMPv3 Access View Configuration).

#### TLV 53 -- SNMPv1v2c Coexistence Sub-TLVs

| Type | CANN Name | Spec |
|------|-----------|------|
| 53.1 | SNMPv1v2c Community Name | DOCSIS 2.0 · MULPIv4.0 I11 |
| 53.2 | SNMPv1v2c Transport Address Access | DOCSIS 2.0 · MULPIv4.0 I11 |
| 53.2.1 | SNMPv1v2c Transport Address | DOCSIS 2.0 · MULPIv4.0 I11 |
| 53.2.2 | SNMPv1v2c Transport Address Mask | DOCSIS 2.0 · MULPIv4.0 I11 |
| 53.3 | SNMPv1v2c Access View Type | DOCSIS 2.0 · MULPIv4.0 I11 |
| 53.4 | SNMPv1v2c Access View Name | DOCSIS 2.0 · MULPIv4.0 I11 |

#### TLV 54 -- SNMPv3 Access View Configuration Sub-TLVs

| Type | CANN Name | Spec |
|------|-----------|------|
| 54.1 | SNMPv3 Access View Name | DOCSIS 2.0 · MULPIv4.0 I11 |
| 54.2 | SNMPv3 Access View Subtree | DOCSIS 2.0 · MULPIv4.0 I11 |
| 54.3 | SNMPv3 Access View Mask | DOCSIS 2.0 · MULPIv4.0 I11 |
| 54.4 | SNMPv3 Access View Type | DOCSIS 2.0 · MULPIv4.0 I11 |

### 3.10 TLV 72 -- MESP Sub-TLVs

CANN Section 11.1.7. Sub-TLVs within TLV 72 (Metro Ethernet Service Profile).

| Type | CANN Name | Spec |
|------|-----------|------|
| 72.1 | MESP Reference | DPoE-MULPIv2.0 I14 |
| 72.2 | MESP Bandwidth Profile (MESP-BP) | DPoE-MULPIv2.0 I14 |
| 72.2.1 | MESP-BP Committed Information Rate | DPoE-MULPIv2.0 I14 |
| 72.2.2 | MESP-BP Committed Burst Size | DPoE-MULPIv2.0 I14 |
| 72.2.3 | MESP-BP Excess Information Rate | DPoE-MULPIv2.0 I14 |
| 72.2.4 | MESP-BP Excess Burst Size | DPoE-MULPIv2.0 I14 |
| 72.2.5 | MESP-BP Coupling Flag | DPoE-MULPIv2.0 I14 |
| 72.2.6 | MESP-BP Color Mode | DPoE-MULPIv2.0 I14 |
| 72.2.6.1 | MESP-BP-CM Color Identification Field | DPoE-MULPIv2.0 I14 |
| 72.2.6.2 | MESP-BP-CM Color Identification Field Value | DPoE-MULPIv2.0 I14 |
| 72.2.7 | MESP-BP Color Marking | DPoE-MULPIv2.0 I14 |
| 72.2.7.1 | MESP-BP-CR Color Marking Field | DPoE-MULPIv2.0 I14 |
| 72.2.7.2 | MESP-BP-CR Color Marking Field Value | DPoE-MULPIv2.0 I14 |
| 72.3 | MESP Name | DPoE-MULPIv2.0 I14 |

### 3.11 TLV 83 -- L2CP Sub-TLVs

CANN Section 11.1.8. Sub-TLVs within TLV 83 (L2CP Management).

| Type | CANN Name | Spec |
|------|-----------|------|
| 83.1 | CMIM | DPoE-MULPIv2.0 I14 |
| 83.2 | L2CP Mode | DPoE-MULPIv2.0 I14 |
| 83.3 | L2CP L2PT D-MAC Address | DPoE-MULPIv2.0 I14 |
| 83.4 | L2CP Filter | DPoE-MULPIv2.0 I14 |

### 3.12 TLV 202 -- eRouter Sub-TLVs

CANN Section 11.1.9. Sub-TLVs within TLV 202 (eRouter).

| Type | CANN Name | Spec |
|------|-----------|------|
| 202.1 | eRouter Initialization Mode Encoding | eRouter I22 |
| 202.2 | TR-069 Management Server | eRouter I22 |
| 202.2.1 | EnableCWMP | eRouter I22 |
| 202.2.2 | URL Parameter | eRouter I22 |
| 202.2.3 | Username Parameter | eRouter I22 |
| 202.2.4 | Password Parameter | eRouter I22 |
| 202.2.5 | Connection Request Username | eRouter I22 |
| 202.2.6 | Connection Request Password | eRouter I22 |
| 202.2.7 | ACSOverride | eRouter I22 |
| 202.3 | eRouter Initialization Mode Override | eRouter I22 |
| 202.10 | Router Advertisement (RA) Transmission Interval | eRouter I22 |
| 202.11 | SNMP MIB Object | eRouter I22 |
| 202.12 | IP Multicast Configuration Server | eRouter I22 |
| 202.13 | Link-ID Control | eRouter I22 |
| 202.42 | Topology Mode Encoding | eRouter I22 |
| 202.43 | Vendor Specific Information | eRouter I22 |
| 202.43.8 | Vendor ID Encoding | eRouter I22 |
| 202.53 | SNMPv1v2c Coexistence Configuration | eRouter I22 |
| 202.53.1 | SNMPv1v2c Community Name | eRouter I22 |
| 202.53.2 | SNMPv1v2c Community Name | eRouter I22 |
| 202.53.2.1 | SNMPv1v2c Transport Address | eRouter I22 |
| 202.53.2.2 | SNMPv1v2c Transport Address Mask | eRouter I22 |
| 202.53.2.3 | SNMPv1v2c Access View Type | eRouter I22 |
| 202.53.2.4 | SNMPv1v2c Access View Name | eRouter I22 |
| 202.54 | SNMPv3 Access View Configuration | eRouter I22 |
| 202.54.1 | SNMPv3 Access View Name | eRouter I22 |
| 202.54.2 | SNMPv3 Access View Subtree | eRouter I22 |
| 202.54.3 | SNMPv3 Access View Mask | eRouter I22 |
| 202.54.4 | SNMPv3 Access View Type | eRouter I22 |

### 3.13 TLV 219 -- eTEA Sub-TLVs

CANN Section 11.1.10. Sub-TLVs within TLV 219 (eTEA -- TDM Emulation Adapter). All sub-TLVs are defined in CM-SP-TEI.

| Type | CANN Name | Spec |
|------|-----------|------|
| 219.8 | eTEA Symbol Clock | TEI I06 |
| 219.9 | eTEA IWF cfg-encoding | TEI I06 |
| 219.9.1 | eTEA PW Index Setting | TEI I06 |
| 219.9.2 | eTEA Emulation Type | TEI I06 |
| 219.9.3 | eTEA PW Peer Address | TEI I06 |
| 219.9.4 | eTEA PW Peer Address IPv6 | TEI I06 |
| 219.9.5 | eTEA PW Destination Port | TEI I06 |
| 219.9.6 | eTEA PW Peer Destination Port | TEI I06 |
| 219.9.7 | eTEA PW Name | TEI I06 |
| 219.9.8 | eTEA PW Description | TEI I06 |
| 219.9.9 | eTEA PW Admin Status | TEI I06 |
| 219.9.10 | eTEA Status Change Notification Enable | TEI I06 |
| 219.9.11 | eTEA PW TDM Type | TEI I06 |
| 219.9.12 | eTEA PW TDM Configuration Table Index | TEI I06 |
| 219.9.13 | eTEA PW CESoPSNConfiguration Index | TEI I06 |
| 219.9.14 | eTEA PW RTP SSRC | TEI I06 |
| 219.9.15 | eTEA PW Peer SSRC | TEI I06 |
| 219.9.16 | eTEA PW TDM Circulation Map | TEI I06 |
| 219.9.16.1 | eTEA PW TDM Port # | TEI I06 |
| 219.9.16.2 | eTEA PW TDM Timeslot Map | TEI I06 |
| 219.10 | eTEA PW TDM Configuration Table | TEI I06 |
| 219.10.1 | eTEA PW TDM Configuration Table Index | TEI I06 |
| 219.10.2 | eTEA PW TDM Payload Size | TEI I06 |
| 219.10.3 | eTEA PW RTP Header Used | TEI I06 |
| 219.10.5 | eTEA Jitter Buffer Setting | TEI I06 |
| 219.10.6 | eTEA PW Payload Suppression | TEI I06 |
| 219.10.7 | eTEA PW LOPS Exit Criteria | TEI I06 |
| 219.10.8 | eTEA PW LOPS Entrance Criteria | TEI I06 |
| 219.10.10 | eTEA PW Packet Replace Policy | TEI I06 |
| 219.10.11 | eTEA PW Packet Loss Window | TEI I06 |
| 219.10.12 | eTEA PW Excessive Loss Threshold | TEI I06 |
| 219.10.15 | eTEA PW Severe Loss Threshold | TEI I06 |
| 219.10.16 | eTEA PW RTP Timestamp Mode | TEI I06 |
| 219.10.17 | eTEA PW Default Fill Pattern | TEI I06 |
| 219.10.18 | eTEA PW L Flag Payload Policy | TEI I06 |
| 219.10.19 | eTEA PW TOS | TEI I06 |
| 219.10.20 | eTEA PW RTP Payload Type | TEI I06 |
| 219.10.21 | eTEA PW RTP Peer Payload Type | TEI I06 |
| 219.10.22 | eTEA PW RTP Timestamp Reference | TEI I06 |
| 219.10.23 | eTEA PW RTP Peer Timestamp Reference | TEI I06 |
| 219.10.24 | eTEA PW SRTP Enable | TEI I06 |
| 219.11 | eTEA SNMP MIB Object | TEI I06 |
| 219.12 | eTEA SNMP Write-Access Control | TEI I06 |
| 219.13 | CESoPSN Configuration Table | TEI I06 |
| 219.13.1 | PW CESoPSN Config Table Index | TEI I06 |
| 219.13.2 | Default Idle Pattern | TEI I06 |
| 219.13.3 | LFlag Policy | TEI I06 |
| 219.13.4 | RFlag Policy | TEI I06 |
| 219.13.5 | Remote Defect Policy | TEI I06 |
| 219.13.6 | LOPS Policy | TEI I06 |
| 219.13.7 | App Sig TOS | TEI I06 |
| 219.13.8 | RTP CAS PT | TEI I06 |
| 219.13.9 | RTP CAS Peer PT | TEI I06 |
| 219.13.10 | App Sig Idle | TEI I06 |
| 219.13.11 | App Sig Interval | TEI I06 |
| 219.13.12 | App Sig Max Interval | TEI I06 |
| 219.14 | dsx1 Configuration Table | TEI I06 |
| 219.14.1 | dsx1 Port ID | TEI I06 |
| 219.14.2 | dsx1 Line Type | TEI I06 |
| 219.14.3 | dsx1 Line Coding | TEI I06 |
| 219.14.4 | dsx1 Circuit ID | TEI I06 |
| 219.14.5 | dsx1 Loopback Configuration | TEI I06 |
| 219.14.6 | dsx1 Signal Mode | TEI I06 |
| 219.14.7 | dsx1 Transmit Clock Source | TEI I06 |
| 219.14.8 | dsx1 Fdl | TEI I06 |
| 219.14.9 | dsx1 Line Length | TEI I06 |
| 219.14.10 | dsx1 Line Status Trap Enable | TEI I06 |
| 219.14.11 | dsx1 Channelization | TEI I06 |
| 219.14.12 | dsx1 Line Mode | TEI I06 |
| 219.14.13 | dsx1 Line Build Out | TEI I06 |
| 219.43 | eTEA Vendor Specific Extensions | TEI I06 |
| 219.255 | eTEA End of Text | TEI I06 |

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

| Type | CANN Name | Reference |
|------|-----------|-----------|
| 43.1 | CM Load Balancing Policy ID | C.1.1.18.1.1 |
| 43.2 | CM Load Balancing Priority | C.1.1.18.1.2 |
| 43.3 | CM Load Balancing Group ID | C.1.1.18.1.3 |
| 43.4 | CM Ranging Class ID Extension | C.1.1.18.1.4 |
| 43.5 | L2VPN Encoding | C.1.1.18.1.5 |
| 43.6 | Extended CMTS MIC Configuration Setting | C.1.1.18.1.6 |
| 43.7 | Source Address Verification (SAV) | C.1.1.18.1.7 |
| 43.8 | Cable Modem Attribute Masks | C.1.1.18.1.8 |
| 43.9 | IP Multicast Join Authorization | C.1.1.18.1.9 |
| 43.10 | Service Type Identifier | C.1.1.18.1.10 |
| 43.12 | DEMARC Auto-Configuration (DAC) | C.1.1.18.1.11 |

### Reference

- MULPIv4.0 C.1.1.18 -- DOCSIS Extension Field
- MULPIv4.0 C.1.1.18.1 -- General Extension Information Encodings
- MULPIv4.0 C.1.1.18.2 -- Vendor-Specific Encodings
- CANN Section 11.1.2 -- TLV 43 Sub-TLV Definitions
