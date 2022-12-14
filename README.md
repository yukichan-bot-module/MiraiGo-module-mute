# MiraiGo-module-mute

ID: `com.aimerneige.mute`

Module for [MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 功能

- 通过发送指定消息快速禁言群成员

指令示例

- 开启全员禁言/开启全体禁言
- 关闭全员禁言/关闭全体禁言
- 禁言 **@对象**
- **@对象** 禁言
- 禁言30 **@对象**
- **@对象** 禁言30

> **Note**
>
> 禁言 0 即为取消禁言

## 使用方法

在适当位置引用本包

```go
package example

imports (
    // ...

    _ "github.com/yukichan-bot-module/MiraiGo-module-mute"

    // ...
)

// ...
```

在全局配置文件中修改默认禁言时间（单位分钟，不写默认为 2）

```yaml
aimerneige:
  mute:
    default: 5
```

## LICENSE

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
<img src="https://www.gnu.org/graphics/agplv3-155x51.png">
</a>

本项目使用 `AGPLv3` 协议开源，您可以在 [GitHub](https://github.com/yukichan-bot-module/MiraiGo-module-mute) 获取本项目源代码。为了整个社区的良性发展，我们强烈建议您做到以下几点：

- **间接接触（包括但不限于使用 `Http API` 或 跨进程技术）到本项目的软件使用 `AGPLv3` 开源**
- **不鼓励，不支持一切商业使用**
