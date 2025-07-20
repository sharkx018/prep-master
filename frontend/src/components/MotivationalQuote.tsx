import React from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { Quote } from 'lucide-react';
import { useMotivationalQuote } from '../hooks/useMotivationalQuote';

const MotivationalQuote: React.FC = () => {
  const { isDarkMode } = useTheme();
  const selectedQuote = useMotivationalQuote();

  if (!selectedQuote) return null;

  // Split the quote to separate the text from the author
  const quoteParts = selectedQuote.split(' – ');
  const quoteText = quoteParts[0];
  const author = quoteParts[1] || '';

  return (
    <div className={`rounded-xl shadow-lg p-6 mb-8 border transition-colors duration-300 ${
      isDarkMode 
        ? 'bg-gradient-to-br from-indigo-900/30 to-purple-900/30 border-indigo-800' 
        : 'bg-gradient-to-br from-indigo-50 to-purple-50 border-indigo-100'
    }`}>
      <div className="flex items-start space-x-4">
        <div className={`p-3 rounded-full shadow-sm flex-shrink-0 ${
          isDarkMode ? 'bg-indigo-800/50' : 'bg-white'
        }`}>
          <Quote className="h-6 w-6 text-indigo-600" />
        </div>
        <div className="flex-1">
          <blockquote className={`text-lg font-medium leading-relaxed mb-3 ${
            isDarkMode ? 'text-gray-100' : 'text-gray-900'
          }`}>
            "{quoteText}"
          </blockquote>
          {author && (
            <cite className={`text-sm font-semibold ${
              isDarkMode ? 'text-indigo-400' : 'text-indigo-600'
            }`}>
              — {author}
            </cite>
          )}
        </div>
      </div>
    </div>
  );
};

export default MotivationalQuote; 