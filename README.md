# BK-GSE-SDK

[![license](https://img.shields.io/badge/license-mit-brightgreen.svg?style=flat)](https://github.com/TencentBlueKing/bk-gse-sdk/blob/master/LICENSE.txt)
[![Release Version](https://img.shields.io/badge/release-v2-brightgreen.svg)](https://github.com/TencentBlueKing/bk-gse-sdk/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TencentBlueKing/bk-gse-sdk/pulls)
[![Blueking Pipelines Status](https://api.bkdevops.qq.com/process/api/external/pipelines/projects/gse/p-61ae051c1fb24c36a1c570c29eec7b28/badge?X-DEVOPS-PROJECT-ID=gse)](https://api.bkdevops.qq.com/process/api-html/user/builds/projects/gse/pipelines/p-61ae051c1fb24c36a1c570c29eec7b28/latestFinished?X-DEVOPS-PROJECT-ID=gse)

[English](README_EN.md) | 简体中文

> **重要提示**: `master` 分支在开发过程中可能处于**不稳定或者不可用状态**。
请通过[releases](https://github.com/TencentBlueKing/bk-gse-sdk/releases) 而非 `master` 去获取稳定的版本。

蓝鲸智云管控平台（BlueKing General Service Engine）简称GSE，本项目提供GSE的各项SDK。

## Overview

* [golang-sdk](go/README.md)

## Features

### 平台化消息槽

- Server端的消息下行投递、上行消息回调
- Agent端的下行消息接收、上行消息投递

如果想了解以上功能的详细说明，请参考[功能说明](https://bk.tencent.com/docs)

## Getting started

* [上下行信令通信](docs/plugin_message.md)

## Roadmap

* [版本日志](CHANGELOG.md)

## Support

- [白皮书](https://bk.tencent.com/docs)
- [社区论坛](https://bk.tencent.com/s-mart/community)

## BlueKing Community

- [BK-JOB](https://github.com/TencentBlueKing/bk-job) 蓝鲸智云作业平台(Job)是一套运维脚本管理系统，具备海量任务并发处理能力。
- [BK-CMDB](https://github.com/TencentBlueKing/bk-cmdb)：蓝鲸智云配置平台（蓝鲸智云 CMDB）是一个面向资产及应用的企业级配置管理平台。
- [BK-CI](https://github.com/TencentBlueKing/bk-ci)：蓝鲸智云持续集成平台是一个开源的持续集成和持续交付系统，可以轻松将你的研发流程呈现到你面前。
- [BK-BCS](https://github.com/TencentBlueKing/bk-bcs)：蓝鲸智云容器管理平台是以容器技术为基础，为微服务业务提供编排管理的基础服务平台。
- [BK-PaaS](https://github.com/TencentBlueKing/blueking-paas)：蓝鲸智云 PaaS 平台是一个开放式的开发平台，让开发者可以方便快捷地创建、开发、部署和管理 SaaS 应用。
- [BK-SOPS](https://github.com/TencentBlueKing/bk-sops)：蓝鲸智云标准运维（SOPS）是通过可视化的图形界面进行任务流程编排和执行的系统，是蓝鲸智云体系中一款轻量级的调度编排类 SaaS 产品。

## Contributing

如果你有好的意见或建议，欢迎给我们提 Issues 或 Pull Requests，为蓝鲸智云开源社区贡献力量。
请阅读[Contributing Guide](CONTRIBUTING.md)了解项目开发基础规范参与代码贡献。

[腾讯开源激励计划](https://opensource.tencent.com/contribution)鼓励开发者的参与和贡献，期待你的加入。

## License

项目基于 MIT 协议，详细请参考 [LICENSE](LICENSE.txt)。

我们承诺未来不会更改适用于交付给任何人的当前项目版本的开源许可证（MIT 协议）。
