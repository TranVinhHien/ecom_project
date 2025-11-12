import createNextIntlPlugin from "next-intl/plugin";
// import * from ""
/** @type {import('next').NextConfig} */
const withNextIntl = createNextIntlPlugin()
const nextConfig = {
    output: "standalone",
};

export default withNextIntl(nextConfig);
