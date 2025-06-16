import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  formatters: {
    css: true,
    html: false,
    markdown: true,
  },
}, [
  {
    rules: {
      '@typescript-eslint/consistent-type-definitions': ['error', 'type'],
      'no-console': 'off',
      'import/first': 'off',
      'no-alert': 'off',
    },
  },
])
