import { jsPDF } from "jspdf";
import autoTable from "jspdf-autotable";

// Helper function to format time
const formatTime = (timeString) => {
  if (!timeString) return "";
  const [hours, minutes] = timeString.split(":");
  return `${parseInt(hours) % 12 || 12}:${minutes}${hours >= 12 ? "pm" : "am"}`;
};

// Generate PDF report
export const generatePdfReport = (
  dayDetails,
  selectedDate,
  selectedCourse,
  selectedFaculty,
  selectedSemester,
  courses,
  faculties,
  semesters
) => {
  // Initialize jsPDF
  const doc = new jsPDF();
  
  // Add autoTable to jsPDF instance
  autoTable(doc, {});

  // Header section
  doc.setFontSize(18);
  doc.text("Attendance Calendar Report", 14, 22);
  
  doc.setFontSize(12);
  const dateString = selectedDate.toLocaleDateString("en-US", {
    weekday: "long",
    month: "long",
    day: "numeric",
    year: "numeric",
  });
  doc.text(`Date: ${dateString}`, 14, 32);
  
  // Filters information
  let filters = [];
  if (selectedCourse !== "all") {
    const courseName = courses.find((c) => c.ID == selectedCourse)?.Name || selectedCourse;
    filters.push(`Course: ${courseName}`);
  }
  if (selectedFaculty !== "all") {
    const facultyName = faculties.find((f) => f.ID == selectedFaculty)?.Name || selectedFaculty;
    filters.push(`Faculty: ${facultyName}`);
  }
  if (selectedSemester !== "all") {
    const semesterName = semesters.find((s) => s.ID == selectedSemester)?.Name || selectedSemester;
    filters.push(`Semester: ${semesterName}`);
  }
  
  if (filters.length > 0) {
    doc.text(`Filters: ${filters.join(", ")}`, 14, 42);
  }
  
  // Class details section
  doc.setFontSize(16);
  doc.text("Class Details", 14, 60);
  
  // Table data - keep raw status for custom rendering
  const tableData = dayDetails.map((detail) => [
    detail.subject || "N/A",
    detail.faculty || "N/A",
    `${formatTime(detail.start_time)} - ${formatTime(detail.end_time)}`,
    detail.course_name || "N/A",
    detail.semester ? `Semester ${detail.semester}` : "N/A",
    detail.room || "N/A",
    `${detail.course_name || ""} ${detail.batch_year || ""} ${detail.batch_section || ""}`.trim() || "N/A",
    detail.status || "" // Raw status value for custom rendering
  ]);
  
  // Table headers
  const headers = [
    "Subject",
    "Faculty",
    "Time",
    "Course",
    "Semester",
    "Room",
    "Batch",
    "Status",
  ];
  
  // Generate table using the autoTable function directly
  autoTable(doc, {
    startY: 65,
    head: [headers],
    body: tableData,
    theme: "grid",
    styles: {
      fontSize: 10,
      cellPadding: 2,
      overflow: 'linebreak',
      halign: 'left',
    },
    headStyles: { fillColor: [52, 73, 94] },
    columnStyles: {
      0: { cellWidth: 30 }, // Subject
      1: { cellWidth: 25 }, // Faculty
      2: { cellWidth: 25 }, // Time
      3: { cellWidth: 20 }, // Course
      4: { cellWidth: 20 }, // Semester
      5: { cellWidth: 15 }, // Room
      6: { cellWidth: 25 }, // Batch
      7: { cellWidth: 20 }, // Status
    },
    didDrawCell: function(data) {
      // Custom rendering for status column (index 7)
      if (data.column.index === 7 && data.cell.section === 'body') {
        const status = data.cell.raw;
        const colors = {
          'held': '#28a745',      // green for Class Taken
          'cancelled': '#dc3545', // red for Class Missed
          '': '#6c757d'           // gray for No Entry
        };
        const color = colors[status] || '#6c757d';
        
        doc.setFillColor(color);
        doc.rect(
  data.cell.x,
  data.cell.y,
  data.cell.width,
  data.cell.height,
  'F'
);

      }
    },
    didDrawPage: function (data) {
      // Footer
      doc.setFontSize(10);
      doc.text(
        `Report generated on ${new Date().toLocaleDateString()}`,
        data.settings.margin.left,
        doc.internal.pageSize.height - 10
      );
    },
  });
  
  // Save and open PDF
  const fileName = `Attendance_Report_${selectedDate.toISOString().split('T')[0]}.pdf`;
  doc.save(fileName);
  
  // Open in new tab
  const pdfBlob = doc.output("blob");
  const pdfUrl = URL.createObjectURL(pdfBlob);
  window.open(pdfUrl, "_blank");
};

// Date Ranged function
export const generateDateRangePdfReport = async (
  startDate,
  endDate,
  selectedCourse,
  selectedFaculty,
  selectedSemester,
  courses,
  faculties,
  semesters
) => {
  try {
    // Fetch data for the entire date range
    const API_BASE_URL = "http://localhost:8080/api/v1";
    const params = new URLSearchParams({
      start_date: startDate.toISOString().split('T')[0],
      end_date: endDate.toISOString().split('T')[0]
    });
    
    if (selectedSemester !== "all") params.append('semester', selectedSemester);
    if (selectedCourse !== "all") params.append('course_id', selectedCourse);
    if (selectedFaculty !== "all") params.append('faculty_id', selectedFaculty);

    const response = await fetch(`${API_BASE_URL}/calendar/range?${params}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include'
    });

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
    const data = await response.json();
    
    // Initialize jsPDF
    const doc = new jsPDF();
    autoTable(doc, {});
    
    // Header section
    doc.setFontSize(18);
    doc.text("Attendance Calendar Report", 14, 22);
    
    doc.setFontSize(12);
    const startStr = startDate.toLocaleDateString("en-US", {
      weekday: "long",
      month: "long",
      day: "numeric",
      year: "numeric",
    });
    const endStr = endDate.toLocaleDateString("en-US", {
      weekday: "long",
      month: "long",
      day: "numeric",
      year: "numeric",
    });
    doc.text(`Date Range: ${startStr} to ${endStr}`, 14, 32);
    
    // Filters information
    let filters = [];
    if (selectedCourse !== "all") {
      const courseName = courses.find((c) => c.ID == selectedCourse)?.Name || selectedCourse;
      filters.push(`Course: ${courseName}`);
    }
    if (selectedFaculty !== "all") {
      const facultyName = faculties.find((f) => f.ID == selectedFaculty)?.Name || selectedFaculty;
      filters.push(`Faculty: ${facultyName}`);
    }
    if (selectedSemester !== "all") {
      const semesterName = semesters.find((s) => s.ID == selectedSemester)?.Name || selectedSemester;
      filters.push(`Semester: ${semesterName}`);
    }
    
    if (filters.length > 0) {
      doc.text(`Filters: ${filters.join(", ")}`, 14, 42);
    }
    
    // Summary section
    doc.setFontSize(16);
    doc.text("Summary", 14, 60);
    
    // Calculate totals
    const totals = data.reduce((acc, day) => {
      acc.held += day.total_held || 0;
      acc.cancelled += day.total_cancelled || 0;
      acc.noData += day.no_data || 0;
      return acc;
    }, { held: 0, cancelled: 0, noData: 0 });
    
    const summaryData = [
      ["Total Classes Held", totals.held],
      ["Total Classes Cancelled", totals.cancelled],
      ["Total No Data", totals.noData],
      ["Total Classes", totals.held + totals.cancelled + totals.noData]
    ];
    
    autoTable(doc, {
      startY: 65,
      head: [['Metric', 'Count']],
      body: summaryData,
      theme: 'grid',
      headStyles: { fillColor: [52, 73, 94] },
      styles: { fontSize: 12, cellPadding: 3 },
      columnStyles: { 0: { cellWidth: 70 }, 1: { cellWidth: 30 } }
    });
    
    // Detailed section
    doc.setFontSize(16);
    doc.text("Daily Details", 14, doc.lastAutoTable.finalY + 15);
    
    // Prepare detailed data with raw status values
    const tableData = data.flatMap(day => {
      return day.details.map(detail => [
        day.date,
        detail.subject || "N/A",
        detail.faculty || "N/A",
        `${formatTime(detail.start_time)} - ${formatTime(detail.end_time)}`,
        detail.course_name || "N/A",
        detail.semester ? `Semester ${detail.semester}` : "N/A",
        detail.room || "N/A",
        `${detail.course_name || ""} ${detail.batch_year || ""} ${detail.batch_section || ""}`.trim() || "N/A",
        detail.status || "" // Raw status value for custom rendering
      ]);
    });
    
    // Table headers
    const headers = [
      "Date",
      "Subject",
      "Faculty",
      "Time",
      "Course",
      "Semester",
      "Room",
      "Batch",
      "Status",
    ];
    
    // Generate table with custom status rendering
    autoTable(doc, {
      startY: doc.lastAutoTable.finalY + 20,
      head: [headers],
      body: tableData,
      theme: "grid",
      styles: {
        fontSize: 10,
        cellPadding: 2,
        overflow: 'linebreak',
        halign: 'left',
      },
      headStyles: { fillColor: [52, 73, 94] },
      columnStyles: {
        0: { cellWidth: 20 }, // Date
        1: { cellWidth: 25 }, // Subject
        2: { cellWidth: 25 }, // Faculty
        3: { cellWidth: 20 }, // Time
        4: { cellWidth: 20 }, // Course
        5: { cellWidth: 15 }, // Semester
        6: { cellWidth: 15 }, // Room
        7: { cellWidth: 25 }, // Batch
        8: { cellWidth: 20 }, // Status
      },
      didDrawCell: function(data) {
        // Custom rendering for status column (index 8)
        if (data.column.index === 8 && data.cell.section === 'body') {
          const status = data.cell.raw;
          const colors = {
            'held': '#28a745',      // green for Class Taken
            'cancelled': '#dc3545', // red for Class Missed
            '': '#6c757d'           // gray for No Entry
          };
          const color = colors[status] || '#6c757d';
          
          doc.setFillColor(color);
          doc.rect(
  data.cell.x,
  data.cell.y,
  data.cell.width,
  data.cell.height,
  'F'
);
        }
      },
      didDrawPage: function (data) {
        // Footer
        doc.setFontSize(10);
        doc.text(
          `Report generated on ${new Date().toLocaleDateString()}`,
          data.settings.margin.left,
          doc.internal.pageSize.height - 10
        );
      },
    });
    
    // Save and open PDF
    const fileName = `Attendance_Report_${startDate.toISOString().split('T')[0]}_to_${endDate.toISOString().split('T')[0]}.pdf`;
    doc.save(fileName);
    
    // Open in new tab
    const pdfBlob = doc.output("blob");
    const pdfUrl = URL.createObjectURL(pdfBlob);
    window.open(pdfUrl, "_blank");
    
  } catch (error) {
    console.error("Error generating date range PDF:", error);
    alert(`Failed to generate PDF: ${error.message}`);
  }
};