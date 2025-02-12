import js from "@eslint/js";
import globals from "globals";
import react from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import checkFile from "eslint-plugin-check-file";
import importPlugin from "eslint-plugin-import";
import prettier from "eslint-plugin-prettier";
import tseslint from "typescript-eslint";

export default tseslint.config({
  settings: { react: { version: "18.3" } },
  extends: [
    js.configs.recommended,
    ...tseslint.configs.recommendedTypeChecked,
    ...tseslint.configs.stylisticTypeChecked,
    prettier,
  ],
  files: ["**/*.{ts,tsx}"],
  languageOptions: {
    ecmaVersion: "latest",
    globals: globals.browser,
    parserOptions: {
      project: "./tsconfig.json",
      tsconfigRootDir: import.meta.dirname,
    },
  },
  plugins: {
    "react-hooks": reactHooks,
    "react-refresh": reactRefresh,
    react,
    "check-file": checkFile,
    import: importPlugin,
    prettier,
  },
  rules: {
    ...reactHooks.configs.recommended.rules,
    "react-refresh/only-export-components": [
      "warn",
      { allowConstantExport: true },
    ],
    ...react.configs.recommended.rules,
    ...react.configs["jsx-runtime"].rules,
    quotes: [
      "error",
      "double",
      {
        avoidEscape: true,
        allowTemplateLiterals: true,
      },
    ],
    "import/no-restricted-paths": [
      "error",
      {
        zones: [
          // disables cross-feature imports:
          // eg. src/features/discussions should not import from src/features/comments, etc.
          {
            target: "./src/features/auth",
            from: "./src/features",
            except: ["./auth"],
          },
          {
            target: "./src/features/comments",
            from: "./src/features",
            except: ["./comments"],
          },
          {
            target: "./src/features/discussions",
            from: "./src/features",
            except: ["./discussions"],
          },
          {
            target: "./src/features/teams",
            from: "./src/features",
            except: ["./teams"],
          },
          {
            target: "./src/features/users",
            from: "./src/features",
            except: ["./users"],
          },
          // unidirectiona codebase: e.g. src/app can import from src/features but not the other way around
          {
            target: "./src/features",
            from: "./src/app",
          },
          // e.g src/features and src/app can import from these shared modules but not the other way around
          {
            target: [
              "./src/components",
              "./src/hooks",
              "./src/lib",
              "./src/types",
              "./src/utils",
            ],
            from: ["./src/features", "./src/app"],
          },
        ],
      },
    ],
    "import/no-cycle": "error",
    "linebreak-style": ["error", "unix"],
    "react/prop-types": "off",
    "import/order": [
      "error",
      {
        groups: [
          "builtin",
          "external",
          "internal",
          "parent",
          "sibling",
          "index",
          "object",
        ],
        "newlines-between": "always",
        alphabetize: { order: "asc", caseInsensitive: true },
      },
    ],
    "import/default": "off",
    "import/no-named-as-default-member": "off",
    "import/no-named-as-default": "off",
    "react/react-in-jsx-scope": "off",
    "jsx-a11y/anchor-is-valid": "off",
    "@typescript-eslint/no-unused-vars": ["error"],
    "@typescript-eslint/explicit-function-return-type": ["off"],
    "@typescript-eslint/explicit-module-boundary-types": ["off"],
    "@typescript-eslint/no-empty-function": ["off"],
    "@typescript-eslint/no-explicit-any": ["off"],
    "prettier/prettier": ["error", {}, { usePrettierrc: true }],
    "check-file/filename-naming-convention": [
      "error",
      {
        "**/*.{ts,tsx}": "KEBAB_CASE",
      },
      {
        ignoreMiddleExtensions: true,
      },
    ],
    "check-file/filename-naming-convention": [
      "error",
      {
        "**/*.{ts,tsx}": "KEBAB_CASE",
      },
      {
        // ignore the middle extensions of the filename to support filename like bable.config.js or smoke.spec.ts
        ignoreMiddleExtensions: true,
      },
    ],
    "check-file/folder-naming-convention": [
      "error",
      {
        // all folders within src (except __tests__)should be named in kebab-case
        "src/**/!(__tests__)": "KEBAB_CASE",
      },
    ],
  },
});
