!url: piconani
!create: 2024-10-22

#### https://pic.onani.cn

#### 是一个由二叉树树搭建的二次元随机图网站，网站提供了横屏竖屏图API获取，图源在Cloudflare R2

1. 最早，这个网站搭建在一台VPS上，图源在本地服务器，通过PHP来调用。之后，由于带宽和资费问题，逐步迈向新架构

2. 随后，将图源迁移到了Github仓库，通过Cloudlflare CDN代理访问，随之而来的问题是，整个仓库高达500M，并且每一次更新图片都需要进行很久的构建，速度也不算快
   
   - 这期间，图源从 `.png` 转换为了 `.webp` ，最终大小200M左右，骤降60%

3. 最终，采用Cloudflare R2作为图源（免费10G），连接Cloudflare Workers，得到两个域名（ `hrandom.onani.cn` 、 `vrandom.onani.cn` ），再通过创建一个Cloudflare Pages完成整个网站

目前，这个随机图网站仅需一个域名即可存活。~~我嘎了网站还活着~~

### 花絮

1. 在使用PHP的阶段时，尝试购买多种大带宽服务器而浪费了近百元

2. 在某一个阶段，尝试使用AList将各大网盘抽象为对象存储作为图源，最终因为仍然需要一个本地服务器和涉嫌滥用最终不了了之

3. 不久之前，曾尝试部署到IPFS，然后发现是[灵车](/pages-html.html#fleek等基于ipfs的托管商-灵车)
