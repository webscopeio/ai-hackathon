import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      keyframes: {
        glow: {
          "0%, 100%": {
            borderColor: "rgb(59 130 246 / 0.7)",
            boxShadow: "0 0 15px rgba(59,130,246,0.5)",
          },
          "50%": {
            borderColor: "rgb(59 130 246 / 0.4)",
            boxShadow: "0 0 15px rgba(59,130,246,0.3)",
          },
        },
      },
      animation: {
        glow: "glow 3s ease-in-out infinite",
      },
    },
  },
  plugins: [],
};

export default config;
