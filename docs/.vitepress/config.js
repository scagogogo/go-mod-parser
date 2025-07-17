import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Go Mod Parser',
  description: 'A comprehensive Go module parser library',
  
  // Base URL for GitHub Pages
  base: '/go-mod-parser/',
  
  // Language configuration
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'Go Mod Parser',
      description: 'A comprehensive Go module parser library',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'API Reference', link: '/api/' },
          { text: 'Examples', link: '/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/go-mod-parser' }
        ],
        sidebar: [
          {
            text: 'Getting Started',
            items: [
              { text: 'Introduction', link: '/' },
              { text: 'Installation', link: '/installation' },
              { text: 'Quick Start', link: '/quick-start' }
            ]
          },
          {
            text: 'API Reference',
            items: [
              { text: 'Overview', link: '/api/' },
              { text: 'Core Functions', link: '/api/core-functions' },
              { text: 'Data Structures', link: '/api/data-structures' },
              { text: 'Helper Functions', link: '/api/helper-functions' },
              { text: 'Error Handling', link: '/api/error-handling' }
            ]
          },
          {
            text: 'Examples',
            items: [
              { text: 'Overview', link: '/examples/' },
              { text: 'Basic Parsing', link: '/examples/basic-parsing' },
              { text: 'File Discovery', link: '/examples/file-discovery' },
              { text: 'Dependency Analysis', link: '/examples/dependency-analysis' },
              { text: 'Advanced Usage', link: '/examples/advanced-usage' }
            ]
          }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2023 Software Composition Analysis'
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'Go Mod Parser',
      description: '全面的 Go 模块解析库',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: 'API 参考', link: '/zh/api/' },
          { text: '示例', link: '/zh/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/go-mod-parser' }
        ],
        sidebar: [
          {
            text: '开始使用',
            items: [
              { text: '介绍', link: '/zh/' },
              { text: '安装', link: '/zh/installation' },
              { text: '快速开始', link: '/zh/quick-start' }
            ]
          },
          {
            text: 'API 参考',
            items: [
              { text: '概览', link: '/zh/api/' },
              { text: '核心函数', link: '/zh/api/core-functions' },
              { text: '数据结构', link: '/zh/api/data-structures' },
              { text: '辅助函数', link: '/zh/api/helper-functions' },
              { text: '错误处理', link: '/zh/api/error-handling' }
            ]
          },
          {
            text: '示例',
            items: [
              { text: '概览', link: '/zh/examples/' },
              { text: '基础解析', link: '/zh/examples/basic-parsing' },
              { text: '文件发现', link: '/zh/examples/file-discovery' },
              { text: '依赖分析', link: '/zh/examples/dependency-analysis' },
              { text: '高级用法', link: '/zh/examples/advanced-usage' }
            ]
          }
        ],
        footer: {
          message: '基于 MIT 许可证发布。',
          copyright: '版权所有 © 2023 Software Composition Analysis'
        }
      }
    }
  },
  
  themeConfig: {
    search: {
      provider: 'local'
    }
  }
})
