import React, { useState, useEffect } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { 
  Trophy, 
  Clock, 
  ExternalLink, 
  Minimize2, 
  Maximize2, 
  AlertCircle,
  Calendar,
  Timer,
  RefreshCw,
  Bell,
  BellOff
} from 'lucide-react';
import { leetcodeApi, LeetCodeContest } from '../services/api';

interface FloatingContestWidgetProps {}

const CountdownTimer: React.FC<{ targetTimestamp: number; compact?: boolean }> = ({ targetTimestamp, compact = false }) => {
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
        setIsLive(true);
        return { days: 0, hours: 0, minutes: 0, seconds: 0 };
      }
      
      return {
        days: Math.floor(difference / (1000 * 60 * 60 * 24)),
        hours: Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
        minutes: Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60)),
        seconds: Math.floor((difference % (1000 * 60)) / 1000)
      };
    };

    setTimeLeft(calculateTimeLeft());
    
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft());
    }, 1000);
    
    return () => clearInterval(timer);
  }, [targetTimestamp]);

  if (isLive) {
    return (
      <div className="flex items-center">
        <div className="flex items-center bg-gradient-to-r from-green-500 to-emerald-500 text-white px-2 py-1 rounded-full shadow-lg animate-pulse">
          <div className="w-2 h-2 bg-white rounded-full mr-2 animate-ping"></div>
          <span className={`${compact ? 'text-xs' : 'text-sm'} font-bold`}>LIVE NOW!</span>
        </div>
      </div>
    );
  }

  const formatTime = () => {
    if (compact) {
      if (timeLeft.days > 0) return `${timeLeft.days}d ${timeLeft.hours}h`;
      if (timeLeft.hours > 0) return `${timeLeft.hours}h ${timeLeft.minutes}m`;
      return `${timeLeft.minutes}m ${timeLeft.seconds}s`;
    }
    
    return `${timeLeft.days > 0 ? `${timeLeft.days}d ` : ''}${timeLeft.hours > 0 ? `${timeLeft.hours}h ` : ''}${timeLeft.minutes}m ${timeLeft.seconds}s`;
  };

  return (
    <div className="flex items-center">
      <div className="flex items-center bg-gradient-to-r from-indigo-500 to-purple-500 text-white px-2 py-1 rounded-full shadow-md">
        <Clock className={`${compact ? 'h-3 w-3' : 'h-4 w-4'} mr-1`} />
        <span className={`${compact ? 'text-xs' : 'text-sm'} font-medium`}>
          {formatTime()}
        </span>
      </div>
    </div>
  );
};

const FloatingContestWidget: React.FC<FloatingContestWidgetProps> = () => {
  const { isDarkMode } = useTheme();
  const [contests, setContests] = useState<LeetCodeContest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isMinimized, setIsMinimized] = useState(false);
  const [notifications, setNotifications] = useState(true);
  const [lastFetch, setLastFetch] = useState(Date.now());

  const fetchContests = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await leetcodeApi.getContests();
      console.log('Fetched contests:', data); // Debug log
      setContests(data);
      setLastFetch(Date.now());
    } catch (err) {
      setError('Failed to fetch contests');
      console.error('Error fetching contests:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContests();
    
    // Auto-refresh every 5 minutes
    const interval = setInterval(fetchContests, 5 * 60 * 1000);
    
    return () => clearInterval(interval);
  }, []);

  const isUpcoming = (startTime: number) => {
    return startTime * 1000 > Date.now();
  };

  const isStartingSoon = (startTime: number) => {
    const timeDiff = startTime * 1000 - Date.now();
    return timeDiff > 0 && timeDiff <= 24 * 60 * 60 * 1000; // Within 24 hours
  };

  const upcomingContests = contests
    .filter(contest => isUpcoming(contest.startTime))
    .slice(0, 3); // Show only next 3 contests
  
  // Debug log
  console.log('Upcoming contests:', upcomingContests);

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  };

  const formatTimeToIST = (timestamp: number, short = false) => {
    const date = new Date(timestamp * 1000);
    
    if (short) {
      const options: Intl.DateTimeFormatOptions = {
        day: 'numeric',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit',
        timeZone: 'Asia/Kolkata'
      };
      return date.toLocaleString('en-IN', options);
    }
    
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

  // Request notification permission
  const requestNotificationPermission = async () => {
    if ('Notification' in window && Notification.permission === 'default') {
      await Notification.requestPermission();
    }
  };

  useEffect(() => {
    if (notifications) {
      requestNotificationPermission();
    }
  }, [notifications]);

  // Show notification for contests starting soon
  useEffect(() => {
    if (!notifications || !('Notification' in window) || Notification.permission !== 'granted') {
      return;
    }

    upcomingContests.forEach(contest => {
      if (isStartingSoon(contest.startTime)) {
        const timeDiff = contest.startTime * 1000 - Date.now();
        const hoursLeft = Math.floor(timeDiff / (1000 * 60 * 60));
        
        if (hoursLeft <= 1 && hoursLeft > 0) {
          new Notification(`Contest Starting Soon!`, {
            body: `${contest.title} starts in ${hoursLeft} hour(s)`,
            icon: '/favicon.ico',
            tag: contest.titleSlug
          });
        }
      }
    });
  }, [upcomingContests, notifications]);

  if (isMinimized) {
    return (
      <div className={`fixed bottom-4 right-4 z-[9999] backdrop-blur-md ${
        isDarkMode 
          ? 'bg-gray-900/90 border-gray-600' 
          : 'bg-white/90 border-gray-200'
      } border rounded-2xl shadow-2xl p-4 max-w-xs transform transition-all duration-300 hover:scale-105 hover:shadow-3xl`} 
      style={{ zIndex: 9999 }}>
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <div className="flex items-center mr-4">
              <div className="w-2 h-2 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full animate-pulse mr-2"></div>
              <Trophy className="h-5 w-5 text-yellow-500" />
            </div>
            <div>
              <span className={`text-sm font-semibold ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
                {upcomingContests.length} Contest{upcomingContests.length !== 1 ? 's' : ''}
              </span>
              {upcomingContests.length > 0 && (
                <div className="text-xs text-gray-500">Next starting soon</div>
              )}
            </div>
          </div>
          <button
            onClick={() => setIsMinimized(false)}
            className={`p-2 rounded-full transition-all duration-200 hover:scale-110 ${
              isDarkMode 
                ? 'hover:bg-gray-700 text-gray-300 hover:text-white' 
                : 'hover:bg-gray-100 text-gray-600 hover:text-gray-900'
            }`}
            title="Expand"
          >
            <Maximize2 className="h-4 w-4" />
          </button>
        </div>
        
        {upcomingContests.length > 0 && (
          <div className="mt-3 pt-2 border-t border-dashed border-gray-300">
            <CountdownTimer targetTimestamp={upcomingContests[0].startTime} compact={true} />
          </div>
        )}
      </div>
    );
  }

  return (
    <div className={`fixed bottom-4 right-4 z-[9999] backdrop-blur-xl ${
      isDarkMode 
        ? 'bg-gradient-to-br from-gray-900/95 via-gray-800/95 to-gray-900/95 border-gray-600' 
        : 'bg-gradient-to-br from-white/95 via-gray-50/95 to-white/95 border-gray-200'
    } border rounded-3xl shadow-2xl max-w-sm w-full max-h-96 overflow-hidden transform transition-all duration-300 hover:shadow-3xl hover:scale-[1.02]`} 
    style={{ zIndex: 9999 }}>
      {/* Header */}
      <div className={`p-4 border-b ${isDarkMode ? 'border-gray-600/50' : 'border-gray-200/50'} flex items-center justify-between bg-gradient-to-r ${isDarkMode ? 'from-gray-800/50 to-gray-700/50' : 'from-gray-50/50 to-white/50'}`}>
        <div className="flex items-center">
          <div className="flex items-center mr-4">
            <div className="w-2 h-2 bg-gradient-to-r from-yellow-400 to-orange-500 rounded-full animate-pulse mr-2"></div>
            <Trophy className="h-6 w-6 text-yellow-500" />
          </div>
          <div>
            <h3 className={`font-bold text-lg ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
              Live Contests
            </h3>
            <p className={`text-xs ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
              Track upcoming contests
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-1">
          <button
            onClick={() => setNotifications(!notifications)}
            className={`p-2 rounded-full transition-all duration-200 hover:scale-110 ${
              isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-gray-100'
            }`}
            title={notifications ? 'Disable notifications' : 'Enable notifications'}
          >
            {notifications ? (
              <Bell className="h-4 w-4 text-indigo-500" />
            ) : (
              <BellOff className="h-4 w-4 text-gray-400" />
            )}
          </button>
          <button
            onClick={fetchContests}
            className={`p-2 rounded-full transition-all duration-200 hover:scale-110 ${
              isDarkMode ? 'hover:bg-gray-700' : 'hover:bg-gray-100'
            }`}
            title="Refresh"
            disabled={loading}
          >
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin text-indigo-500' : 'text-gray-500'}`} />
          </button>
          <button
            onClick={() => setIsMinimized(true)}
            className={`p-2 rounded-full transition-all duration-200 hover:scale-110 ${
              isDarkMode ? 'hover:bg-gray-700 text-gray-400' : 'hover:bg-gray-100 text-gray-500'
            }`}
            title="Minimize"
          >
            <Minimize2 className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="max-h-80 overflow-y-auto">
        {loading ? (
          <div className="p-6 flex flex-col items-center justify-center">
            <div className="relative">
              <RefreshCw className="h-8 w-8 animate-spin text-indigo-500" />
              <div className="absolute inset-0 h-8 w-8 animate-ping bg-indigo-500/20 rounded-full"></div>
            </div>
            <span className={`text-sm mt-3 font-medium ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
              Loading contests...
            </span>
          </div>
        ) : error ? (
          <div className="p-6 text-center">
            <div className={`p-3 rounded-full inline-flex mb-3 ${isDarkMode ? 'bg-red-900/20' : 'bg-red-50'}`}>
              <AlertCircle className={`h-6 w-6 ${isDarkMode ? 'text-red-400' : 'text-red-500'}`} />
            </div>
            <div className={`text-sm font-medium mb-3 ${isDarkMode ? 'text-red-300' : 'text-red-600'}`}>
              {error}
            </div>
            <button
              onClick={fetchContests}
              className="inline-flex items-center px-3 py-2 text-sm font-medium text-white bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full hover:from-indigo-600 hover:to-purple-600 transition-all duration-200 transform hover:scale-105"
            >
              <RefreshCw className="h-4 w-4 mr-1" />
              Try again
            </button>
          </div>
        ) : upcomingContests.length === 0 ? (
          <div className="p-6 text-center">
            <div className={`p-4 rounded-full inline-flex mb-4 ${isDarkMode ? 'bg-gray-700/50' : 'bg-gray-100'}`}>
              <Calendar className={`h-8 w-8 ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`} />
            </div>
            <p className={`text-sm font-medium ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
              No upcoming contests
            </p>
            <p className={`text-xs mt-1 ${isDarkMode ? 'text-gray-500' : 'text-gray-400'}`}>
              Check back later for new contests
            </p>
          </div>
        ) : (
          <div className="p-4 space-y-4">
            {upcomingContests.map((contest, index) => (
              <div
                key={contest.titleSlug}
                className={`relative p-4 rounded-2xl border transition-all duration-300 transform hover:scale-[1.02] hover:shadow-lg ${
                  isDarkMode
                    ? 'bg-gradient-to-br from-gray-800/50 to-gray-700/50 border-gray-600/50 hover:border-indigo-400/50'
                    : 'bg-gradient-to-br from-white/80 to-gray-50/80 border-gray-200/50 hover:border-indigo-300/50'
                } ${index === 0 ? 'ring-2 ring-indigo-500/20' : ''}`}
              >
                {index === 0 && (
                  <div className="absolute -top-2 -right-2">
                    <div className="bg-gradient-to-r from-indigo-500 to-purple-500 text-white text-xs font-bold px-3 py-1 rounded-full shadow-lg animate-bounce">
                      Next Up!
                    </div>
                  </div>
                )}
                
                <div className="mb-4">
                  <h4 className={`text-base font-bold ${isDarkMode ? 'text-white' : 'text-gray-900'} leading-tight mb-1`}>
                    <a
                      href={`https://leetcode.com/contest/${contest.titleSlug}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:text-indigo-500 transition-colors duration-200 hover:underline block"
                    >
                      {contest.title || contest.titleSlug?.replace(/-/g, ' ').replace(/\b\w/g, l => l.toUpperCase()) || 'LeetCode Contest'}
                    </a>
                  </h4>
                  <p className={`text-xs ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                    LeetCode Programming Contest
                  </p>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className={`flex items-center text-xs ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                      <div className={`p-1 rounded-md mr-2 ${isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}`}>
                        <Calendar className="h-3 w-3" />
                      </div>
                      <span className="font-medium">{formatTimeToIST(contest.startTime, true)}</span>
                    </div>
                    
                    <div className={`flex items-center text-xs ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                      <div className={`p-1 rounded-md mr-2 ${isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}`}>
                        <Timer className="h-3 w-3" />
                      </div>
                      <span className="font-medium">{formatDuration(contest.duration)}</span>
                    </div>
                  </div>

                  <div className="flex justify-center">
                    <CountdownTimer targetTimestamp={contest.startTime} compact={true} />
                  </div>
                </div>

                <div className="mt-4 pt-3 border-t border-dashed border-gray-300/50">
                  <a
                    href={`https://leetcode.com/contest/${contest.titleSlug}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="w-full inline-flex items-center justify-center px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-indigo-500 to-purple-500 rounded-xl hover:from-indigo-600 hover:to-purple-600 transition-all duration-200 transform hover:scale-105 shadow-md hover:shadow-lg"
                  >
                    <Trophy className="h-3 w-3 mr-2" />
                    Join Contest
                    <ExternalLink className="h-3 w-3 ml-2" />
                  </a>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className={`px-4 py-3 border-t ${isDarkMode ? 'border-gray-600/50 bg-gradient-to-r from-gray-800/30 to-gray-700/30' : 'border-gray-200/50 bg-gradient-to-r from-gray-50/30 to-white/30'} text-center`}>
        <div className="flex items-center justify-center space-x-2">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
          <span className={`text-xs font-medium ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
            Updated {new Date(lastFetch).toLocaleTimeString()}
          </span>
        </div>
      </div>
    </div>
  );
};

export default FloatingContestWidget;
