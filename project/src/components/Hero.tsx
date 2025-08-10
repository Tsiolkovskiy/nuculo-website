import React from 'react';
import { ArrowRight, Shield, Zap } from 'lucide-react';

export default function Hero() {
  return (
    <div className="relative bg-gradient-to-br from-gray-50 to-blue-50 pt-32 pb-20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            Smart & Secure Software
            <span className="text-blue-600"> Solutions</span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Empowering businesses with cutting-edge technology solutions that drive growth and innovation.
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-4">
            <button className="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors">
              Get Started
              <ArrowRight className="ml-2" size={20} />
            </button>
            <button className="inline-flex items-center px-6 py-3 border-2 border-blue-600 text-blue-600 rounded-md hover:bg-blue-50 transition-colors">
              Learn More
            </button>
          </div>
        </div>

        <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="bg-white p-6 rounded-lg shadow-md">
            <Shield className="text-blue-600 mb-4" size={32} />
            <h3 className="text-xl font-semibold mb-2">Secure by Design</h3>
            <p className="text-gray-600">Built-in security at every layer of your software infrastructure.</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-md">
            <Zap className="text-blue-600 mb-4" size={32} />
            <h3 className="text-xl font-semibold mb-2">High Performance</h3>
            <p className="text-gray-600">Optimized solutions that deliver exceptional speed and reliability.</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-md">
            <Shield className="text-blue-600 mb-4" size={32} />
            <h3 className="text-xl font-semibold mb-2">Enterprise Ready</h3>
            <p className="text-gray-600">Scalable solutions designed for growing businesses.</p>
          </div>
        </div>
      </div>
    </div>
  );
}