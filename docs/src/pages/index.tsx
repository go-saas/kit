import React from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import HomepageFeatures from "@site/src/components/HomepageFeatures";

import styles from "./index.module.css";

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <header className={clsx("hero hero--primary", styles.heroBanner)}>
      <div className="container">
        <img src={require('@site/static/img/logo.png').default} style={{height:50,width:50}} />
        <h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className={`button button--secondary button--lg ${styles.buttonsItems}`}
            to="/docs/intro/overview"
          >
            <div>Tutorial</div>
          </Link>
          <Link
             className={`button button--secondary button--lg ${styles.buttonsItems}`}
            to="https://saas.nihaosaoya.com/"
          >
            <div>Demo</div>
          </Link>
        </div>
       
      </div>
    </header>
  );
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`[GO-SAAS-KIT] ${siteConfig.title}`}
      description="Full-featured open source SAAS starter kit"
    >
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
