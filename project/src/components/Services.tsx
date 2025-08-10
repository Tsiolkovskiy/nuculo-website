import React from 'react';
import { Code, Database, Cloud, Globe, Shield, Settings } from 'lucide-react';

export default function Services() {
  const services = [
    {
      icon: <Code className="text-blue-600" size={32} />,
      title: 'Custom Software Development',
      description: 'Tailored solutions built to address your unique business challenges.'
    },
    {
      icon: <Cloud className="text-blue-600" size={32} />,
      title: 'Cloud Solutions',
      description: 'Scalable cloud infrastructure and migration services.'
    },
    {
      icon: <Database className="text-blue-600" size={32} />,
      title: 'Data Analytics',
      description: 'Transform your data into actionable insights.'
    },
    {
      icon: <Globe className="text-blue-600" size={32} />,
      title: 'Web Applications',
      description: 'Modern web applications built with cutting-edge technologies.'
    },
    {
      icon: <Shield className="text-blue-600" size={32} />,
      title: 'Cybersecurity',
      description: 'Comprehensive security solutions to protect your digital assets.'
    },
    {
      icon: <Settings className="text-blue-600" size={32} />,
      title: 'DevOps Services',
      description: 'Streamline your development and deployment processes.'
    }
  ];

  return (
    <section id="services" className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-bold text-gray-900 mb-4">Our Services</h2>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Comprehensive software solutions tailored to your business needs
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {services.map((service, index) => (
            <div key={index} className="p-6 border border-gray-200 rounded-lg hover:shadow-lg transition-shadow">
              <div className="mb-4">{service.icon}</div>
              <h3 className="text-xl font-semibold mb-2">{service.title}</h3>
              <p className="text-gray-600">{service.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}