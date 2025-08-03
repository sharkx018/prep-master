import { useState, useEffect } from 'react';

// Move quotes array to a separate constant
const QUOTES = [
  // Albert Einstein
  "Genius is 1% talent and 99% hard work. – Albert Einstein",
  "It's not that I'm so smart, it's just that I stay with problems longer. – Albert Einstein",
  // "Strive not to be a success, but rather to be of value. – Albert Einstein",

  "You will have to hang tough thru this period, but all lies ahead of you, donot stop until you got answers to all question and make your life a success story! - Ashish Vishnoi",

  "Har pehlu zindagi ka intemn hota h, muskilo se bachna aasan hota h, Dar kr bhadne walo ko kuch nhi milta zindagi me, Ladne walo k kadmo m jagah hota h. -- some where in hostel room(jee)",

  
  "It’s not that I’m so smart, it’s just that I stay with problems longer. – Albert Einstein",
  // "In the middle of difficulty lies opportunity. – Albert Einstein",
  // "Imagination is more important than knowledge. For knowledge is limited, whereas imagination embraces the entire world. – Albert Einstein",
  // "A person who never made a mistake never tried anything new. – Albert Einstein",
  // "The only source of knowledge is experience. – Albert Einstein",
  // "Try not to become a man of success, but rather try to become a man of value. – Albert Einstein",
  // "Weakness of attitude becomes weakness of character. – Albert Einstein",
  // "Life is like riding a bicycle. To keep your balance, you must keep moving. – Albert Einstein",
  // "Education is what remains after one has forgotten what one has learned in school. – Albert Einstein",
  // "Only those who attempt the absurd can achieve the impossible. – Albert Einstein",
  // "You never fail until you stop trying. – Albert Einstein",
  // "The important thing is not to stop questioning. Curiosity has its own reason for existing. – Albert Einstein",
  "I have no special talent. I am only passionately curious. – Albert Einstein",

  // Thomas Edison
  // "Opportunity is missed by most people because it is dressed in overalls and looks like work. – Thomas Edison",
  // "I have not failed. I've just found 10,000 ways that won't work. – Thomas Edison",
  "There is no substitute for hard work. – Thomas Edison",
  "Hard work beats Talent. - xyz",

  // Bruce Lee
  "I fear not the man who has practiced 10,000 kicks once, but I fear the man who has practiced one kick 10,000 times. – Bruce Lee",
  "The successful warrior is the average man, with laser-like focus. – Bruce Lee",
  
  // Others
  "Work like hell. I mean you just have to put in 80 to 100 hour weeks every week. – Elon Musk",
  "Working 16 hrs a day, 7 Days a week, 52 Weeks in the year and People still call me lucky",

];

const STORAGE_KEY = 'motivationalQuote';

// Global variable to store the quote in memory as well
let globalQuote: string | null = null;

export const useMotivationalQuote = () => {
  const [quote, setQuote] = useState<string>('');

  useEffect(() => {
    // If we already have a quote in memory, use it
    if (globalQuote) {
      setQuote(globalQuote);
      return;
    }

    // Check sessionStorage
    const storedQuote = sessionStorage.getItem(STORAGE_KEY);
    
    if (storedQuote && QUOTES.includes(storedQuote)) {
      globalQuote = storedQuote;
      setQuote(storedQuote);
    } else {
      // Generate new quote
      const randomIndex = Math.floor(Math.random() * QUOTES.length);
      const newQuote = QUOTES[randomIndex];
      
      // Store in both sessionStorage and global variable
      sessionStorage.setItem(STORAGE_KEY, newQuote);
      globalQuote = newQuote;
      setQuote(newQuote);
    }
  }, []);

  return quote;
}; 