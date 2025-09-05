import React, { useState, useEffect } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { Clock, Calendar, ExternalLink, Trophy, Timer, RefreshCw, ChevronLeft, ChevronRight, AlertCircle } from 'lucide-react';
import { leetcodeApi, LeetCodeContest } from '../services/api';
import MotivationalQuote from '../components/MotivationalQuote';

// Countdown timer component
const CountdownTimer: React.FC<{ targetTimestamp: number }> = ({ targetTimestamp }) => {
  const [timeLeft, setTimeLeft] = useState<{ days: number; hours: number; minutes: number; seconds: number }>({
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0
  });
  const [isLive, setIsLive] = useState(false);

  useEffect(() => {
    const calculateTimeLeft = () => {
      const now = new Date().getTime();
      const targetTime = targetTimestamp * 1000;
      const difference = targetTime - now;
      
      if (difference <= 0) {
        // Contest has started
        setIsLive(true);
        return { days: 0, hours: 0, minutes: 0, seconds: 0 };
      }
      
      // Calculate time left
      return {
        days: Math.floor(difference / (1000 * 60 * 60 * 24)),
        hours: Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
        minutes: Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60)),
        seconds: Math.floor((difference % (1000 * 60)) / 1000)
      };
    };

    // Initial calculation
    setTimeLeft(calculateTimeLeft());
    
    // Update every second
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft());
    }, 1000);
    
    return () => clearInterval(timer);
  }, [targetTimestamp]);

  if (isLive) {
    return (
      <div className="flex items-center text-green-500">
        <AlertCircle className="h-5 w-5 mr-2" />
        <span className="text-base font-bold">LIVE NOW!</span>
      </div>
    );
  }

  return (
    <div className="flex items-center">
      <span className="text-base text-indigo-400 font-medium">
        {timeLeft.days > 0 && `${timeLeft.days}d `}
        {timeLeft.hours > 0 && `${timeLeft.hours}h `}
        {timeLeft.minutes}m {timeLeft.seconds}s
      </span>
    </div>
  );
};

const Contest: React.FC = () => {
  const { isDarkMode } = useTheme();
  const [contests, setContests] = useState<LeetCodeContest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  
  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const pastContestsPerPage = 10;

  const fetchContests = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await leetcodeApi.getContests();
      setContests(data);
    } catch (err) {
      setError('Failed to fetch contest data. Please try again later.');
      console.error('Error fetching contests:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContests();
  }, []);

  const formatTimeToIST = (timestamp: number) => {
    // Create date in UTC, then format it for IST timezone
    const date = new Date(timestamp * 1000);
    
    // Format to display in a more readable format: "Sun, Aug 3, 08:00 IST"
    const options: Intl.DateTimeFormatOptions = {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'Asia/Kolkata'
    };
    
    return date.toLocaleString('en-IN', options) + ' IST';
  };

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  };

  const isUpcoming = (startTime: number) => {
    return startTime * 1000 > Date.now();
  };

  const upcomingContests = contests
    .filter(contest => isUpcoming(contest.startTime))
    .sort((a, b) => a.startTime - b.startTime); // Sort upcoming contests by start time (ascending - earliest first)
  
  const allPastContests = contests
    .filter(contest => !isUpcoming(contest.startTime))
    .sort((a, b) => b.startTime - a.startTime); // Sort past contests by start time (descending - most recent first)
  
  // Calculate pagination
  const totalPages = Math.ceil(allPastContests.length / pastContestsPerPage);
  const indexOfLastContest = currentPage * pastContestsPerPage;
  const indexOfFirstContest = indexOfLastContest - pastContestsPerPage;
  const currentPastContests = allPastContests.slice(indexOfFirstContest, indexOfLastContest);
  
  // Handle page change
  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber);
    
    // Scroll to the top of the past contests section
    const pastContestsSection = document.getElementById('past-contests');
    if (pastContestsSection) {
      pastContestsSection.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div>
      {/* Motivational Quote */}
      <MotivationalQuote />

      <div className="mb-8">
        <h1 className={`text-3xl font-bold ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
          <span className="bg-gradient-to-r from-indigo-500 to-purple-600 bg-clip-text text-transparent">
            LeetCode Contests
          </span>
        </h1>
        <p className={`mt-2 ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
          Stay updated with upcoming and past LeetCode coding contests
        </p>
      </div>

      {loading ? (
        <div className="flex justify-center items-center py-16">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
        </div>
      ) : error ? (
        <div className={`p-8 rounded-lg flex flex-col items-center justify-center ${isDarkMode ? 'bg-red-900/20 text-red-300' : 'bg-red-50 text-red-700'}`}>
          <p className="mb-4 text-center">{error}</p>
          <button
            onClick={fetchContests}
            className="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            <RefreshCw className="h-4 w-4 mr-2" />
            Retry
          </button>
        </div>
      ) : (
        <>
          {/* Upcoming Contests */}
          <div className="mb-10">
            <h2 className={`text-2xl font-bold mb-4 flex items-center ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>
              <Trophy className="mr-2 text-yellow-500" />
              Upcoming Contests
            </h2>

            {upcomingContests.length === 0 ? (
              <p className={`p-4 rounded-lg ${isDarkMode ? 'bg-gray-800 text-gray-300' : 'bg-gray-100 text-gray-600'}`}>
                No upcoming contests at the moment. Check back later!
              </p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {upcomingContests.map((contest) => (
                  <div
                    key={contest.titleSlug}
                    className={`rounded-lg border p-6 transition-all duration-300 transform hover:translate-y-[-5px] hover:shadow-lg ${
                      isDarkMode
                        ? 'bg-gray-800 border-gray-700 hover:border-indigo-500'
                        : 'bg-white border-gray-200 hover:border-indigo-500'
                    }`}
                  >
                    <div className="flex items-start justify-between">
                      <h3 className={`text-xl font-bold mb-3 ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
                        <a 
                          href={`https://leetcode.com/contest/${contest.titleSlug}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="hover:text-indigo-500 transition-colors duration-200"
                        >
                          {contest.title}
                        </a>
                      </h3>
                      <div className="bg-gradient-to-r from-indigo-500 to-purple-600 text-white text-sm font-bold px-2 py-1 rounded-full">
                        UPCOMING
                      </div>
                    </div>
                    
                    <div className="space-y-4 mt-4">
                      <div className={`flex items-center ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                        <Calendar className="h-5 w-5 mr-2 text-indigo-400" />
                        <span className="text-base">{formatTimeToIST(contest.startTime)}</span>
                      </div>
                      
                      <div className={`flex items-center ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                        <Timer className="h-5 w-5 mr-2 text-indigo-400" />
                        <span className="text-base">Duration: {formatDuration(contest.duration)}</span>
                      </div>
                      
                      <div className={`flex items-center ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                        <Clock className="h-5 w-5 mr-2 text-indigo-400" />
                        <span className="text-base mr-1">Starts in:</span>
                        <CountdownTimer targetTimestamp={contest.startTime} />
                      </div>
                    </div>
                    
                    <div className="mt-5 pt-4 border-t border-dashed border-gray-600">
                      <a
                        href={`https://leetcode.com/contest/${contest.titleSlug}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center justify-center w-full py-2 px-4 rounded-lg bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-medium hover:from-indigo-700 hover:to-purple-700 transition-all duration-200"
                      >
                        View Contest <ExternalLink className="h-4 w-4 ml-2" />
                      </a>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Past Contests */}
          <div id="past-contests">
            <div className="flex justify-between items-center mb-4">
              <h2 className={`text-2xl font-bold flex items-center ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>
                <Clock className="mr-2 text-gray-400" />
                Past Contests
              </h2>
              
              {totalPages > 1 && (
                <div className={`text-sm ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                  Showing {indexOfFirstContest + 1}-{Math.min(indexOfLastContest, allPastContests.length)} of {allPastContests.length} contests
                </div>
              )}
            </div>

            <div className={`rounded-lg border overflow-hidden ${
              isDarkMode ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'
            }`}>
              <ul className="divide-y divide-gray-700">
                {currentPastContests.map((contest) => (
                  <li key={contest.titleSlug} className={`p-4 hover:bg-gray-700/20 transition-colors duration-200`}>
                    <div className="flex flex-wrap justify-between items-center">
                      <div className="mr-4">
                        <h3 className={`font-medium ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
                          <a 
                            href={`https://leetcode.com/contest/${contest.titleSlug}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="hover:text-indigo-500 transition-colors duration-200"
                          >
                            {contest.title}
                          </a>
                        </h3>
                        <div className={`mt-1 text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                          {formatTimeToIST(contest.startTime)} · Duration: {formatDuration(contest.duration)}
                        </div>
                      </div>
                      
                      <a
                        href={`https://leetcode.com/contest/${contest.titleSlug}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="mt-2 sm:mt-0 inline-flex items-center text-sm font-medium text-indigo-500 hover:text-indigo-400 transition-colors"
                      >
                        View Results <ExternalLink className="h-3 w-3 ml-1" />
                      </a>
                    </div>
                  </li>
                ))}
              </ul>
              
              {/* Pagination controls */}
              {totalPages > 1 && (
                <div className={`flex justify-center items-center p-4 border-t ${
                  isDarkMode ? 'border-gray-700' : 'border-gray-200'
                }`}>
                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => handlePageChange(currentPage - 1)}
                      disabled={currentPage === 1}
                      className={`p-2 rounded-lg ${
                        isDarkMode 
                          ? 'hover:bg-gray-700 disabled:text-gray-600' 
                          : 'hover:bg-gray-100 disabled:text-gray-300'
                      } disabled:cursor-not-allowed`}
                      aria-label="Previous page"
                    >
                      <ChevronLeft className="h-5 w-5" />
                    </button>
                    
                    <div className="flex space-x-1">
                      {Array.from({ length: totalPages }, (_, i) => i + 1)
                        .filter(page => {
                          // Show current page, first and last page, and pages close to current page
                          return (
                            page === 1 || 
                            page === totalPages || 
                            Math.abs(page - currentPage) <= 1
                          );
                        })
                        .map((page, i, filteredPages) => (
                          <React.Fragment key={page}>
                            {i > 0 && filteredPages[i - 1] !== page - 1 && (
                              <span className={`px-2 py-1 ${isDarkMode ? 'text-gray-500' : 'text-gray-400'}`}>...</span>
                            )}
                            <button
                              onClick={() => handlePageChange(page)}
                              className={`w-8 h-8 flex items-center justify-center rounded-lg transition-colors ${
                                currentPage === page
                                  ? isDarkMode
                                    ? 'bg-indigo-600 text-white'
                                    : 'bg-indigo-100 text-indigo-600'
                                  : isDarkMode
                                    ? 'text-gray-300 hover:bg-gray-700'
                                    : 'text-gray-700 hover:bg-gray-100'
                              }`}
                            >
                              {page}
                            </button>
                          </React.Fragment>
                        ))
                      }
                    </div>
                    
                    <button
                      onClick={() => handlePageChange(currentPage + 1)}
                      disabled={currentPage === totalPages}
                      className={`p-2 rounded-lg ${
                        isDarkMode 
                          ? 'hover:bg-gray-700 disabled:text-gray-600' 
                          : 'hover:bg-gray-100 disabled:text-gray-300'
                      } disabled:cursor-not-allowed`}
                      aria-label="Next page"
                    >
                      <ChevronRight className="h-5 w-5" />
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default Contest; 