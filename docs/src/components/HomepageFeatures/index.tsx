import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Src: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Modularity',
    Src: require('@site/static/img/modularity.png').default,
    description: (
      <>
        Every module is layered designed. You can composite them into services easily
      </>
    ),
  },
  {
    title: 'Microservice',
    Src: require('@site/static/img/microservice.png').default,
    description: (
      <>
        Microservice compatible. Integrate with distributed event, distributed transaction. Containerized and provide observability
      </>
    ),
  },
  {
    title: 'Multi-Tenancy',
    Src: require('@site/static/img/multi-tenancy.png').default,
    description: (
      <>
        Multi-tenancy natively, support different database architectures
      </>
    ),
  },
  {
    title: 'Realtime',
    Src: require('@site/static/img/realtime.png').default,
    description: (
      <>
        Realtime with websocket
      </>
    ),
  },
  {
    title: 'Localization',
    Src: require('@site/static/img/localization.png').default,
    description: (
      <>
        Localize your apps and design multilingual contents
      </>
    ),
  },
  {
    title: 'Background Job',
    Src: require('@site/static/img/backgroundjob.png').default,
    description: (
      <>
        Dispatch background jobs
      </>
    ),
  },
  {
    title: 'Business Module',
    Src: require('@site/static/img/business-module.png').default,
    description: (
      <>
        Integrated with Stripe. Design SAAS plans and products with different price models
      </>
    ),
  },
  {
    title: 'Admin UI',
    Src: require('@site/static/img/admin.png').default,
    description: (
      <>
        Admin UI developed with React and Antd. Manage users, tenant etc.
      </>
    ),
  },
  {
    title: 'Customizable',
    Src: require('@site/static/img/customizable.png').default,
    description: (
      <>
        Extend existing modules to fit your own business. Provide templates and tools to develop your own modules
      </>
    ),
  },
];

function Feature({title, Src, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <img className={styles.featureSvg}  src={Src}/>
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
