const Configuration = {
  /*
   * Resolve and load @commitlint/config-conventional from node_modules.
   */
  extends: ['@commitlint/config-conventional'],
  /*
   * Any rules defined here will override rules from @commitlint/config-conventional
   */
  rules: {
    'body-max-line-length': [0],
  },
}
