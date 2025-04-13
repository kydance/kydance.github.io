# 使用 FFmpeg 转换哔哩哔哩缓存视频为 MP4 格式


{{< admonition type=abstract title="导语" open=true >}}
**本文介绍如何使用 FFmpeg 工具将哔哩哔哩缓存的 m4s 视频文件转换为标准 mp4 格式**
{{< /admonition >}}

<!--more-->

## ffmpeg 安装

[**FFmpeg**](https://github.com/FFmpeg/FFmpeg) is a collection of libraries and tools to process multimedia content such as audio, video, subtitles and related metadata.

[下载](https://ffmpeg.org/download.html) `FFmpeg` 并安装

```bash
# Mac
brew install ffmpeg

# Arch Linux
sudo pacman -S ffmpeg

# Debian
sudo apt install ffmpeg

# Verify
ffmpeg -version
```

## BiliBili 视频缓存

1. 哔哩哔哩缓存的视频文件夹，然而由于文件夹的名字使用的纯数字，无法进行很好地辨认，因此，一般需要使用`修改时间`进行排序。

2. 在已下载的文件夹中，视频文件和音频文件被分割为两个文件，且以 `m4s` 作为后缀名：其中包含 `30280` 字样的一般是音频文件，包含 `100035` 字样的一般是视频文件

3. 由于哔哩哔哩对文件进行了简单加密，因此只需**将文件开头9个0删除即可得到原始文件**

## 使用 FFmpeg 进行视频转码

```bash
# Convert video to mp4
ffmpeg -i video-100035.m4s -i audio-30280.m4s -codec copy output.mp4
```

转化后得到的 `output.mp4` 文件大小与原来两个原始 `m4s` 文件大小之和基本保持一致。

## Reference

- [FFmpeg](https://github.com/FFmpeg/FFmpeg)
- [how-to-install-ffmpeg](https://www.hostinger.com/tutorials/how-to-install-ffmpeg)


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/033530b/  

