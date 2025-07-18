import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import academicData from "../assets/academicData.json";
import { RefreshCcw } from "lucide-react";
import clsx from 'clsx';
import SearchableSelect from "@/components/SearchableSelect";

const ClassTimeTable = () => {
  // State for fetched data from APIs (courses, batches, etc.)
  const [courses, setCourses] = useState([]);
  const [batches, setBatches] = useState([]);
  const [semesters, setSemesters] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [faculties, setFaculties] = useState([]);
  const [rooms, setRooms] = useState([]);

  // State for the currently displayed lectures (timetable data)
  const [lectures, setLectures] = useState([]);
  const [gridData, setGridData] = useState({});
  const [allTimeSlots, setAllTimeSlots] = useState([]);

  // State for the selected filter criteria by the user
  const [selectedFilters, setSelectedFilters] = useState({
    course: null,
    batch: null,
    semester: null,
    faculty: null,
    room: null
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [isLoadingLectures, setIsLoadingLectures] = useState(false);

  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
  const API_ENDPOINTS = {
    GET_COURSE: `${API_BASE_URL}/course`,
    GET_BATCH: `${API_BASE_URL}/batch`,
    GET_SEMESTER: `${API_BASE_URL}/semester`,
    GET_SUBJECT: `${API_BASE_URL}/subject`,
    GET_FACULTY: `${API_BASE_URL}/faculty`,
    GET_ROOM: `${API_BASE_URL}/room`,
    LECTURE: `${API_BASE_URL}/lecture`,
    LECTURE_QUERY: `${API_BASE_URL}/lecture/query`,
  };

  useEffect(() => {
    fetchAllData();
  }, []);

  const groupConsecutiveTimeSlots = (lectures, days, timeSlots) => {
    const groupedData = {};
    days.forEach(day => {
      let currentGroup = null;
      timeSlots.forEach((time, timeIndex) => {
        const [startTime, endTime] = time.split('-');
        const lecture = lectures.find(lec =>
          lec.DayOfWeek === day &&
          lec.StartTime === startTime &&
          lec.EndTime === endTime
        );
        if (lecture) {
          const subjectName = subjects.find(sub => sub.ID === lecture.SubjectID)?.Name || 'N/A';
          const facultyName = faculties.find(fac => fac.ID === lecture.FacultyID)?.Name || 'N/A';
          const groupKey = `${day}-${subjectName}-${facultyName}`;
          if (currentGroup && currentGroup.groupKey === groupKey && currentGroup.endIndex === timeIndex - 1) {
            currentGroup.timeSlots.push(time);
            currentGroup.endIndex = timeIndex;
            groupedData[`${day}-${time}`] = currentGroup;
          } else {
            currentGroup = {
              ...lecture,
              groupKey,
              subject: subjectName,
              faculty: facultyName,
              code: subjects.find(sub => sub.ID === lecture.SubjectID)?.Code || 'N/A',
              room: rooms.find(room => room.ID === lecture.RoomID)?.Name || 'N/A',
              timeSlots: [time],
              startIndex: timeIndex,
              endIndex: timeIndex,
              isGrouped: true
            };
            groupedData[`${day}-${time}`] = currentGroup;
          }
        } else {
          currentGroup = null;
        }
      });
    });
    return groupedData;
  };

  const convertLecturesToGridData = (lectures) => {
    const grid = {};
    const uniqueTimeSlots = new Set(academicData.timeSlots);
    lectures.forEach(lecture => {
      const timeSlot = `${lecture.StartTime}-${lecture.EndTime}`;
      uniqueTimeSlots.add(timeSlot);
      const key = `${lecture.DayOfWeek}-${timeSlot}`;
      grid[key] = {
        id: lecture.ID,
        subject: subjects.find(sub => sub.ID === lecture.SubjectID)?.Name || 'N/A',
        code: subjects.find(sub => sub.ID === lecture.SubjectID)?.Code || 'N/A',
        faculty: faculties.find(fac => fac.ID === lecture.FacultyID)?.Name || 'N/A',
        room: rooms.find(room => room.ID === lecture.RoomID)?.Name || 'N/A',
        startTime: lecture.StartTime,
        endTime: lecture.EndTime
      };
    });
    const sortedTimeSlots = sortTimeSlots([...uniqueTimeSlots]);
    setAllTimeSlots(sortedTimeSlots);
    return grid;
  };

  const parseTimeToMinutes = (timeStr) => {
    if (!timeStr) return 0;
    const [hours, minutes] = timeStr.split(':').map(Number);
    return hours * 60 + minutes;
  };

  const sortTimeSlots = (slots) => {
    if (!Array.isArray(slots)) return [];
    return slots.slice().sort((a, b) => {
      const [startA] = a.split('-');
      const [startB] = b.split('-');
      return parseTimeToMinutes(startA) - parseTimeToMinutes(startB);
    });
  };

  const fetchAllData = async () => {
    setLoading(true);
    setError(null);
    try {
      await Promise.all([
        fetchCourses(),
        fetchBatches(),
        fetchSubjects(),
        fetchFaculties(),
        fetchRooms()
      ]);
      setSemesters(Array.isArray(academicData.semesters) ? academicData.semesters : []);
      setAllTimeSlots(Array.isArray(academicData.timeSlots) ? sortTimeSlots(academicData.timeSlots) : []);
    } catch (err) {
      setError('Failed to fetch initial data');
      console.error('Error fetching initial data:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchData = async (endpoint, setter, errorMsg) => {
    try {
      const response = await fetch(endpoint, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include'
      });
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      setter(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error(`Error fetching ${errorMsg}:`, error);
      setError(`Failed to fetch ${errorMsg}`);
      setter([]);
    }
  };

  const fetchCourses = () => fetchData(API_ENDPOINTS.GET_COURSE, setCourses, 'courses');
  const fetchBatches = () => fetchData(API_ENDPOINTS.GET_BATCH, setBatches, 'batches');
  const fetchSubjects = () => fetchData(API_ENDPOINTS.GET_SUBJECT, setSubjects, 'subjects');
  const fetchFaculties = () => fetchData(API_ENDPOINTS.GET_FACULTY, setFaculties, 'faculties');
  const fetchRooms = () => fetchData(API_ENDPOINTS.GET_ROOM, setRooms, 'rooms');

  const fetchLectures = async () => {
    setIsLoadingLectures(true);
    setError(null);
    setLectures([]);
    setGridData({});

    try {
      const queryParams = new URLSearchParams();
      if (selectedFilters.faculty) {
        queryParams.append('faculty_id', selectedFilters.faculty);
      } else if (selectedFilters.room) {
        queryParams.append('room_id', selectedFilters.room);
      } else if (selectedFilters.batch && selectedFilters.semester) {
        queryParams.append('batch_id', selectedFilters.batch);
        queryParams.append('semester', selectedFilters.semester);
      }

      const response = await fetch(`${API_ENDPOINTS.LECTURE_QUERY}?${queryParams}`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include'
      });

      if (!response.ok) {
        if (response.status === 404) {
          setLectures([]);
        } else {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
      } else {
        const filteredLectures = await response.json();
        setLectures(filteredLectures);
        const gridData = convertLecturesToGridData(filteredLectures);
        setGridData(gridData);
      }
    } catch (err) {
      setError('Failed to fetch lectures');
      console.error('Error fetching lectures:', err);
    } finally {
      setIsLoadingLectures(false);
    }
  };

  const handleCourseChange = (value) => {
    setSelectedFilters({ course: value === "all" ? null : value, batch: null, semester: null, faculty: null, room: null });
  };
  const handleBatchChange = (value) => {
    setSelectedFilters(prev => ({ ...prev, batch: value === "all" ? null : value, semester: null, faculty: null, room: null }));
  };
  const handleSemesterChange = (value) => {
    setSelectedFilters(prev => ({ ...prev, semester: value === "all" ? null : value, faculty: null, room: null }));
  };
  const handleFacultyChange = (value) => {
    setSelectedFilters({ faculty: value === "all" ? null : value, course: null, batch: null, semester: null, room: null });
  };
  const handleRoomChange = (value) => {
    setSelectedFilters({ room: value === "all" ? null : value, course: null, batch: null, semester: null, faculty: null });
  };

  const handleGenerateTimetable = () => {
    if ((selectedFilters.course && selectedFilters.batch && selectedFilters.semester) || selectedFilters.faculty || selectedFilters.room) {
      fetchLectures();
    } else {
      alert("Please select either Faculty, Room, or complete Course+Batch+Semester combination to generate timetable.");
    }
  };

  const handleReset = () => {
    setSelectedFilters({ course: null, batch: null, semester: null, faculty: null, room: null });
    setLectures([]);
    setGridData({});
    setAllTimeSlots(Array.isArray(academicData.timeSlots) ? sortTimeSlots(academicData.timeSlots) : []);
  };

  const days = Array.isArray(academicData.days) ? academicData.days : [];

  const courseOptions = [...courses.map(c => ({ value: c.ID, label: c.Name }))];
  const batchOptions = [...batches.filter(batch => selectedFilters.course ? batch.CourseID === selectedFilters.course : true).map(b => ({ value: b.ID, label: `Batch ${b.Year} - Section ${b.Section}` }))];
  const semesterOptions = [...semesters.map(s => ({ value: s.id || s.number, label: s.name || `Semester ${s.number}` }))];
  const facultyOptions = [...faculties.map(f => ({ value: f.ID, label: f.Name }))];
  const roomOptions = [...rooms.map(r => ({ value: r.ID, label: r.Name }))];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50">
      {error && (
        <div className="flex justify-center items-center py-8">
          <div className="bg-white rounded-xl shadow-lg p-6 text-center max-w-md mx-4">
            <div className="mb-4">
              <div className="mx-auto w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
                <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" /></svg>
              </div>
            </div>
            <p className="text-red-600 font-medium mb-4">{error}</p>
            <Button onClick={fetchAllData} className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">Retry</Button>
          </div>
        </div>
      )}

      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="bg-white rounded-2xl shadow-xl border border-gray-100">
          <div className="bg-gradient-to-r from-indigo-600 to-blue-600 px-6 py-4 flex justify-between items-center rounded-t-xl">
            <div>
              <h1 className="text-xl font-bold text-white">View Timetable</h1>
              <p className="text-indigo-100 text-sm mt-1">View academic schedules</p>
            </div>
          </div>

          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="space-y-2 z-[60]">
                <label className="block text-sm font-semibold text-gray-700">Course</label>
                <SearchableSelect options={courseOptions} value={selectedFilters.course} onSelect={handleCourseChange} placeholder="Select course" />
              </div>
              <div className="space-y-2 z-50">
                <label className="block text-sm font-semibold text-gray-700">Batch</label>
                <SearchableSelect options={batchOptions} value={selectedFilters.batch} onSelect={handleBatchChange} placeholder="Select batch" disabled={!selectedFilters.course || loading} />
              </div>
              <div className="space-y-2 z-40">
                <label className="block text-sm font-semibold text-gray-700">Semester</label>
                <SearchableSelect options={semesterOptions} value={selectedFilters.semester} onSelect={handleSemesterChange} placeholder="Select semester" disabled={!selectedFilters.batch} />
              </div>
              <div className="space-y-2 z-30">
                <label className="block text-sm font-semibold text-gray-700">Faculty</label>
                <SearchableSelect options={facultyOptions} value={selectedFilters.faculty} onSelect={handleFacultyChange} placeholder="Select faculty" disabled={loading} />
              </div>
              <div className="space-y-2 z-20">
                <label className="block text-sm font-semibold text-gray-700">Room</label>
                <SearchableSelect options={roomOptions} value={selectedFilters.room} onSelect={handleRoomChange} placeholder="Select room" disabled={loading} />
              </div>
            </div>

            <div className="mt-8 flex flex-col sm:flex-row gap-4">
              <Button
                className={clsx(
                  "w-full sm:w-auto h-12 font-semibold rounded-xl shadow-lg transition-all duration-300",
                  {
                    "bg-gradient-to-r from-indigo-600 to-blue-600 hover:from-indigo-700 hover:to-blue-700 text-white transform hover:scale-[1.02]":
                      (selectedFilters.course && selectedFilters.batch && selectedFilters.semester) || selectedFilters.faculty || selectedFilters.room,
                    "bg-gray-200 text-gray-500 cursor-not-allowed":
                      !((selectedFilters.course && selectedFilters.batch && selectedFilters.semester) || selectedFilters.faculty || selectedFilters.room)
                  }
                )}
                onClick={handleGenerateTimetable}
                disabled={!((selectedFilters.course && selectedFilters.batch && selectedFilters.semester) || selectedFilters.faculty || selectedFilters.room)}
              >
                Generate Timetable
              </Button>
              <Button
                className="w-full sm:w-auto h-12 font-semibold rounded-xl shadow-lg transition-all duration-300 bg-gray-300 hover:bg-gray-400 text-gray-800 flex items-center justify-center gap-2"
                onClick={handleReset}
                disabled={loading || isLoadingLectures}
              >
                <RefreshCcw size={18} />
                Reset
              </Button>
            </div>
          </div>
        </div>
      </div>

      {(isLoadingLectures || loading) && (
        <div className="flex justify-center items-center py-8">
          <div className="bg-white rounded-xl shadow-lg p-6 flex items-center space-x-4">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
            <p className="text-gray-700 font-medium">
              {loading ? 'Loading initial data...' : 'Loading timetable...'}
            </p>
          </div>
        </div>
      )}

      {lectures.length > 0 && !isLoadingLectures && (
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 pb-8">
          <div className="bg-white rounded-2xl shadow-xl border border-gray-100 overflow-hidden">
            <div className="bg-gradient-to-r from-indigo-600 to-blue-600 px-6 py-4">
              <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                  <h2 className="text-xl font-bold text-white">Timetable</h2>
                  <p className="text-indigo-100 text-sm">View your class schedule</p>
                </div>
              </div>
            </div>

            <div className="lg:hidden">
              <div className="p-4 space-y-4">
                {days.map((day) => {
                  const dayLectures = lectures.filter(lec => lec.DayOfWeek === day);
                  return (
                    <div key={day} className="bg-gray-50 rounded-xl p-4">
                      <h3 className="font-bold text-indigo-700 mb-3 text-lg">{day}</h3>
                      <div className="space-y-2">
                        {dayLectures.map((lecture) => {
                          const timeSlot = `${lecture.StartTime}-${lecture.EndTime}`;
                          const subjectName = subjects.find(sub => sub.ID === lecture.SubjectID)?.Name || 'N/A';
                          const facultyName = faculties.find(fac => fac.ID === lecture.FacultyID)?.Name || 'N/A';
                          const roomName = rooms.find(room => room.ID === lecture.RoomID)?.Name || 'N/A';
                          const subjectCode = subjects.find(sub => sub.ID === lecture.SubjectID)?.Code || 'N/A';
                          return (
                            <div key={`${day}-${lecture.ID}`} className="bg-white rounded-lg p-3 border-2 border-gray-200">
                              <div className="flex justify-between items-start">
                                <div className="flex-1">
                                  <div className="text-sm font-semibold text-gray-600 mb-1">{timeSlot}</div>
                                  <div>
                                    <div className="font-semibold text-indigo-700 text-sm mb-1">{subjectName}</div>
                                    <div className="text-xs text-gray-600 mb-1">{subjectCode}</div>
                                    <div className="text-xs text-gray-500">{facultyName}</div>
                                    <div className="text-xs text-gray-500">{roomName}</div>
                                  </div>
                                </div>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>

            <div className="hidden lg:block p-6">
              <div className="overflow-x-auto rounded-xl border-2 border-gray-200">
                <table className="w-full border-collapse">
                  <thead>
                    <tr className="bg-gradient-to-r from-indigo-50 to-blue-50">
                      <th className="border-r border-gray-200 p-4 text-center font-bold text-indigo-700 bg-white">Day / Time</th>
                      {allTimeSlots.map((time, index) => (
                        <th key={index} className="border-r border-gray-200 p-3 text-center font-semibold text-indigo-700 relative min-w-[140px]">
                          <span className="text-sm font-medium">{time}</span>
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {days.map((day, dayIndex) => {
                      const groupedLectures = groupConsecutiveTimeSlots(lectures, [day], allTimeSlots);
                      return (
                        <tr key={day} className={dayIndex % 2 === 0 ? "bg-white" : "bg-gray-50"}>
                          <td className="border-r border-gray-200 p-4 font-bold text-indigo-700 bg-gradient-to-r from-indigo-50 to-blue-50 text-center">{day}</td>
                          {allTimeSlots.map((time, timeIndex) => {
                            const cellKey = `${day}-${time}`;
                            const groupedLecture = groupedLectures[cellKey];
                            if (groupedLecture?.isGrouped && groupedLecture.timeSlots[0] !== time) return null;
                            const colSpan = groupedLecture?.isGrouped ? groupedLecture.timeSlots.length : 1;
                            return (
                              <td key={cellKey} className={clsx("border-r border-gray-200 p-3 text-center h-24 min-w-[140px]", groupedLecture?.isGrouped ? "bg-blue-50" : "")} colSpan={colSpan}>
                                {groupedLecture ? (
                                  <div className="space-y-1">
                                    <div className="font-semibold text-indigo-700 text-sm leading-tight">{groupedLecture.subject}</div>
                                    <div className="text-xs text-gray-600 font-medium">{groupedLecture.code}</div>
                                    <div className="text-xs text-gray-500">{groupedLecture.faculty}</div>
                                    <div className="text-xs text-gray-500">{groupedLecture.room}</div>
                                    {colSpan > 1 && (<div className="text-xs text-gray-400 mt-1">{time.split('-')[0]} to {groupedLecture.timeSlots[groupedLecture.timeSlots.length - 1].split('-')[1]}</div>)}
                                  </div>
                                ) : (
                                  <div className="text-gray-400 text-sm font-medium">No class</div>
                                )}
                              </td>
                            );
                          })}
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      )}

      {!isLoadingLectures && lectures.length === 0 && (
        <div className="flex justify-center items-center py-8">
          <div className="bg-white rounded-xl shadow-lg p-6 text-center max-w-md mx-4">
            {selectedFilters.course || selectedFilters.batch || selectedFilters.semester || selectedFilters.faculty || selectedFilters.room ? (
              <p className="text-gray-700 font-medium">No timetable found for the selected criteria.</p>
            ) : (
              <p className="text-gray-700 font-medium">Please select criteria to view timetable.</p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ClassTimeTable;
