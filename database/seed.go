package database

import (
	"log"
	"time"

	"course-selection-platform/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase checks if the database is empty. If no users exist,
// it wipes dependent tables and populates the initial seed dataset inside a transaction.
func SeedDatabase(db *gorm.DB) error {
	var count int64

	// Check if the users table exists and already contains seeded data
	if db.Migrator().HasTable(&models.Users{}) {
		if err := db.Model(&models.Users{}).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			log.Println("[INFO] Database already contains data. Skipping seed execution.")
			return nil
		}
	}

	log.Println("[INFO] Database is empty. Starting transactional database seeding...")

	// Execute all wipe and seed operations inside a DB transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// ======== 0. Wipe All Existing Data ========
		tables, err := tx.Migrator().GetTables()
		if err != nil {
			return err
		}

		for _, table := range tables {
			if err := tx.Migrator().DropTable(table); err != nil {
				return err
			}
		}

		if err := tx.AutoMigrate(
			&models.Faculty{},
			&models.Users{},
			&models.Subject{},
			&models.SubjectPrerequisite{},
			&models.ChoosableSubjects{},
			&models.Lecturer{},
			&models.SubjectGroup{},
			&models.StudentPriorityList{},
			&models.EnrollmentTable{},
			&models.LecturerRating{},
			&models.SubjectRecommendation{},
			&models.SubjectDifficulty{},
			&models.ActiveRegistration{},
			&models.InactiveRegistration{},
			&models.StudentSubjectAllowance{},
		); err != nil {
			return err
		}

		log.Println("[INFO] Database dropped and re-migrated cleanly.")

		// ======== 1. Seed Faculties (4 Faculties) ========
		faculties := []models.Faculty{
			{
				Name:        "კომპიუტერული მეცნიერება და IT",
				Description: "კომპიუტერული მეცნიერების, პროგრამული ინჟინერიისა და ინფორმაციული ტექნოლოგიების ფაკულტეტი",
			},
			{
				Name:        "ბიზნესის ადმინისტრირება",
				Description: "ეკონომიკის, მენეჯმენტის, ფინანსებისა და მარკეტინგის ფაკულტეტი",
			},
			{
				Name:        "სამართლის ფაკულტეტი",
				Description: "საჯარო, კერძო, სისხლისა და საერთაშორისო სამართლის ფაკულტეტი",
			},
			{
				Name:        "სოციალურ და პოლიტიკურ მეცნიერებათა ფაკულტეტი",
				Description: "საერთაშორისო ურთიერთობების, ფსიქოლოგიის, ჟურნალისტიკისა და სოციოლოგიის ფაკულტეტი",
			},
		}
		for i := range faculties {
			if err := tx.FirstOrCreate(&faculties[i], models.Faculty{Name: faculties[i].Name}).Error; err != nil {
				return err
			}
		}

		// ======== 2. Seed Users ========
		activeStudentStatus := models.Active

		// Extract Faculty IDs as pointers to assign to users safely
		csFacultyID := &faculties[0].ID
		businessFacultyID := &faculties[1].ID
		lawFacultyID := &faculties[2].ID
		socialFacultyID := &faculties[3].ID

		userSeeds := []struct {
			user          models.Users
			plainPassword string
		}{
			{
				user: models.Users{
					Email:         "admin@university.edu",
					Role:          models.RoleAdmin,
					FirstName:     "სისტემური",
					LastName:      "ადმინისტრატორი",
					Gender:        models.Male,
					PersonalNo:    "01001001001",
					PhoneNumber:   "+995555123456",
					PhotoUrl:      "https://example.com/admin.jpg",
					StudentStatus: nil,
					FacultyID:     nil,
				},
				plainPassword: "AdminSecret123!",
			},
			{
				user: models.Users{
					Email:         "g.beridze@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "გიორგი",
					LastName:      "ბერიძე",
					Gender:        models.Male,
					PersonalNo:    "01001001002",
					PhoneNumber:   "+995555654321",
					PhotoUrl:      "https://example.com/student1.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     csFacultyID, // users[1] - CS
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "n.kapanadze@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "ნინო",
					LastName:      "კაპანაძე",
					Gender:        models.Female,
					PersonalNo:    "01001001003",
					PhoneNumber:   "+995555987654",
					PhotoUrl:      "https://example.com/student2.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     lawFacultyID, // users[2] - Law
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "d.maisuradze@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "დავით",
					LastName:      "მაისურაძე",
					Gender:        models.Male,
					PersonalNo:    "01001001004",
					PhoneNumber:   "+995555112233",
					PhotoUrl:      "https://example.com/student3.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     businessFacultyID, // users[3] - Business
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "m.gelashvili@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "მარიამ",
					LastName:      "გელაშვილი",
					Gender:        models.Female,
					PersonalNo:    "01001001005",
					PhoneNumber:   "+995555334455",
					PhotoUrl:      "https://example.com/student4.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     businessFacultyID, // users[4] - Business
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "l.gogichaishvili@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "ლუკა",
					LastName:      "გოგიჩაიშვილი",
					Gender:        models.Male,
					PersonalNo:    "01001001006",
					PhoneNumber:   "+995555667788",
					PhotoUrl:      "https://example.com/student5.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     socialFacultyID, // users[5] - Social Sciences
				},
				plainPassword: "StudentPass123!",
			},
			// --- Managers for testing GetSubjectStudentOverviewService ---
			{
				user: models.Users{
					Email:         "cs.manager@university.edu",
					Role:          models.RoleManager,
					FirstName:     "თამარ",
					LastName:      "აბაშიძე",
					Gender:        models.Female,
					PersonalNo:    "01001001007",
					PhoneNumber:   "+995555778899",
					PhotoUrl:      "https://example.com/manager_cs.jpg",
					StudentStatus: nil,
					FacultyID:     csFacultyID, // users[6] - CS Manager
				},
				plainPassword: "ManagerPass123!",
			},
			{
				user: models.Users{
					Email:         "bus.manager@university.edu",
					Role:          models.RoleManager,
					FirstName:     "ირაკლი",
					LastName:      "ჯავახიშვილი",
					Gender:        models.Male,
					PersonalNo:    "01001001008",
					PhoneNumber:   "+995555889900",
					PhotoUrl:      "https://example.com/manager_bus.jpg",
					StudentStatus: nil,
					FacultyID:     businessFacultyID, // users[7] - Business Manager
				},
				plainPassword: "ManagerPass123!",
			},
			{
				user: models.Users{
					Email:         "unassigned.manager@university.edu",
					Role:          models.RoleManager,
					FirstName:     "ნატო",
					LastName:      "მელიქიშვილი",
					Gender:        models.Female,
					PersonalNo:    "01001001009",
					PhoneNumber:   "+995555990011",
					PhotoUrl:      "https://example.com/manager_none.jpg",
					StudentStatus: nil,
					FacultyID:     nil, // users[8] - Unassigned Manager (tests 400 error)
				},
				plainPassword: "ManagerPass123!",
			},
			// --- Additional CS Students to populate SubjectGroups ---
			{
				user: models.Users{
					Email:         "a.tsiklauri@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "ალექსანდრე",
					LastName:      "წიკლაური",
					Gender:        models.Male,
					PersonalNo:    "01001001010",
					PhoneNumber:   "+995555101010",
					PhotoUrl:      "https://example.com/student6.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     csFacultyID, // users[9] - CS
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "t.khutsishvili@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "თეკლა",
					LastName:      "ხუციშვილი",
					Gender:        models.Female,
					PersonalNo:    "01001001011",
					PhoneNumber:   "+995555111111",
					PhotoUrl:      "https://example.com/student7.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     csFacultyID, // users[10] - CS
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "s.lomidze@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "საბა",
					LastName:      "ლომიძე",
					Gender:        models.Male,
					PersonalNo:    "01001001012",
					PhoneNumber:   "+995555121212",
					PhotoUrl:      "https://example.com/student8.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     csFacultyID, // users[11] - CS
				},
				plainPassword: "StudentPass123!",
			},
			{
				user: models.Users{
					Email:         "e.kvaratskhelia@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "ელენე",
					LastName:      "კვარაცხელია",
					Gender:        models.Female,
					PersonalNo:    "01001001013",
					PhoneNumber:   "+995555131313",
					PhotoUrl:      "https://example.com/student9.jpg",
					StudentStatus: &activeStudentStatus,
					FacultyID:     csFacultyID, // users[12] - CS
				},
				plainPassword: "StudentPass123!",
			},
		}

		var users []models.Users
		for _, seed := range userSeeds {
			hashedBytes, err := bcrypt.GenerateFromPassword([]byte(seed.plainPassword), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			u := seed.user
			u.Password = string(hashedBytes)

			if err := tx.FirstOrCreate(&u, models.Users{Email: u.Email}).Error; err != nil {
				return err
			}
			users = append(users, u)
		}

		// ======== 3. Seed Subjects (20 per Faculty = 80 Total) ========
		type subjectSeed struct {
			name        string
			description string
			credit      int32
			facultyIdx  int
		}

		rawSubjects := []subjectSeed{
			// --- Faculty 0: Computer Science & IT (20 subjects) ---
			{"შესავალი პროგრამირებაში (Python)", "ალგორითმული აზროვნებისა და პროგრამირების საფუძვლები", 6, 0},
			{"ობიექტზე ორიენტირებული პროგრამირება (Java)", "კლასები, მემკვიდრეობითობა, პოლიმორფიზმი და დიზაინ პატერნები", 6, 0},
			{"მონაცემთა სტრუქტურები და ალგორითმები", "ხეები, გრაფები, სორტირება და ალგორითმული სირთულის ანალიზი", 6, 0},
			{"მონაცემთა ბაზების სისტემები (SQL)", "რელაციური მონაცემთა ბაზები, ნორმალიზაცია და SQL მოთხოვნები", 5, 0},
			{"ვებ-პროგრამირება (Frontend - React)", "თანამედროვე მომხმარებლის ინტერფეისების შექმნა React-ით", 5, 0},
			{"ვებ-პროგრამირება (Backend - Go)", "RESTful API-ებისა და მიკროსერვისების აგება Go-ზე", 6, 0},
			{"კომპიუტერული ქსელები", "OSI მოდელი, TCP/IP პროტოკოლები, როუტინგი და ქსელური არქიტექტურა", 5, 0},
			{"ოპერაციული სისტემები (Linux)", "პროცესები, ნაკადები, მეხსიერების მართვა და ფაილური სისტემები", 5, 0},
			{"პროგრამული უზრუნველყოფის ინჟინერია", "Agile, Scrum, ტესტირება და სისტემების პროექტირება", 5, 0},
			{"კიბერუსაფრთხოების საფუძვლები", "ქსელური შეტევების პრევენცია, ავთენტიფიკაცია და უსაფრთხოების აუდიტი", 5, 0},
			{"ხელოვნური ინტელექტი", "ძიების ალგორითმები, ევრისტიკა და ინტელექტუალური აგენტები", 6, 0},
			{"მანქანური სწავლება (Machine Learning)", "რეგრესია, კლასიფიკაცია, კლასტერიზაცია და ნეირონული ქსელები", 6, 0},
			{"ღრუბლოვანი ტექნოლოგიები (AWS & Docker)", "კონტეინერიზაცია, ორკესტრაცია და ღრუბლოვანი ინფრასტრუქტურა", 5, 0},
			{"მობილური აპლიკაციების დეველოპმენტი", "კროს-პლატფორმული მობილური აპლიკაციების შემუშავება", 5, 0},
			{"კომპიუტერული არქიტექტურა", "პროცესორის მუშაობის პრინციპები, ასემბლერი და ციფრული ლოგიკა", 5, 0},
			{"დისკრეტული მათემატიკა", "ლოგიკა, სიმრავლეთა თეორია, კომბინატორიკა და გრაფთა თეორია", 6, 0},
			{"კალკულუსი I (მათემატიკური ანალიზი)", "ზღვრები, წარმოებულები და ინტეგრალები ინჟინერიისთვის", 6, 0},
			{"წრფივი ალგებრა", "მატრიცები, ვექტორული სივრცეები და საკუთარი მნიშვნელობები", 5, 0},
			{"განაწილებული სისტემები", "კონსენსუსის ალგორითმები, რეპლიკაცია და მაღალი ხელმისაწვდომობა", 6, 0},
			{"ადამიანისა და კომპიუტერის ურთიერთქმედება (UI/UX)", "მომხმარებლის გამოცდილების კვლევა და პროტოტიპირება", 4, 0},

			// --- Faculty 1: Business Administration (20 subjects) ---
			{"მიკროეკონომიკის პრინციპები", "მოთხოვნა-მიწოდება, ელასტიურობა და ბაზრის სტრუქტურები", 5, 1},
			{"მაკროეკონომიკის პრინციპები", "მშპ, ინფლაცია, უმუშევრობა და მონეტარული პოლიტიკა", 5, 1},
			{"ფინანსური აღრიცხვა I", "ბუღალტრული ბალანსი, მოგება-ზარალის უწყისი და ფულადი ნაკადები", 6, 1},
			{"მენეჯერული აღრიცხვა", "ხარჯების კალკულაცია, ბიუჯეტირება და შიდა კონტროლი", 5, 1},
			{"მენეჯმენტის საფუძვლები", "დაგეგმვა, ორგანიზება, ლიდერობა და კონტროლი ორგანიზაციებში", 5, 1},
			{"მარკეტინგის საფუძვლები", "ბაზრის სეგმენტაცია, პოზიციონირება და 4P მარკეტინგული მიქსი", 5, 1},
			{"კორპორაციული ფინანსები", "კაპიტალის სტრუქტურა, ფულის დროითი ღირებულება და ინვესტიციები", 6, 1},
			{"ადამიანური რესურსების მართვა (HRM)", "რეკრუტინგი, შეფასების სისტემები და მოტივაციის მექანიზმები", 5, 1},
			{"ოპერაციული მენეჯმენტი", "მიწოდების ჯაჭვი, მარაგების მართვა და ხარისხის კონტროლი", 5, 1},
			{"ბიზნეს სტატისტიკა", "ალბათობის თეორია, ჰიპოთეზების ტესტირება და რეგრესიული ანალიზი", 5, 1},
			{"სტრატეგიული მენეჯმენტი", "კონკურენტული უპირატესობის მოპოვება და გრძელვადიანი დაგეგმვა", 6, 1},
			{"ციფრული მარკეტინგი", "SEO, SEM, სოციალური მედიის მართვა და ანალიტიკა", 5, 1},
			{"საერთაშორისო ბიზნესი", "გლობალური ბაზრები, ექსპორტ-იმპორტი და კულტურული გარემო", 5, 1},
			{"მეწარმეობა და სტარტაპების მართვა", "ბიზნეს მოდელის შექმნა, MVP და ვენჩურული კაპიტალი", 5, 1},
			{"ორგანიზაციული ქცევა", "ჯგუფური დინამიკა, კონფლიქტების მართვა და ორგანიზაციული კულტურა", 4, 1},
			{"ბიზნეს ეთიკა და კორპორაციული პასუხისმგებლობა", "ეთიკური დილემები ბიზნესში და მდგრადი განვითარება", 4, 1},
			{"საბანკო საქმე და ფინანსური ბაზრები", "კომერციული ბანკების ოპერაციები, ფასიანი ქაღალდები და რისკები", 5, 1},
			{"საგადასახადო საქმე", "საქართველოს საგადასახადო კოდექსი და კომპანიების დაბეგვრა", 5, 1},
			{"პროექტების მართვა", "პროექტის სასიცოცხლო ციკლი, რისკები და რესურსების განაწილება", 5, 1},
			{"ბიზნეს კომუნიკაციები და მოლაპარაკებები", "პრეზენტაციის უნარები, საქმიანი მიმოწერა და მოლაპარაკების ტექნიკა", 4, 1},

			// --- Faculty 2: Law (20 subjects) ---
			{"შესავალი სამართალმცოდნეობაში", "სამართლის არსი, წყაროები, სისტემა და ნორმის სტრუქტურა", 6, 2},
			{"საქართველოს კონსტიტუციური სამართალი", "სახელმწიფო მოწყობა, ხელისუფლების დანაწილება და კონსტიტუციური კონტროლი", 6, 2},
			{"სამოქალაქო სამართლის ზოგადი ნაწილი", "პირები, გარიგებები, ვადები და წარმომადგენლობა", 6, 2},
			{"სანივთო სამართალი", "საკუთრება, მფლობელობა, იპოთეკა, გირავნობა და სერვიტუტი", 5, 2},
			{"ვალდებულებითი სამართალი (ზოგადი ნაწილი)", "ვალდებულების წარმოშობა, შესრულება, დარღვევა და შეწყვეტა", 6, 2},
			{"სახელშეკრულებო სამართალი", "ნასყიდობა, ქირავნობა, იჯარა, ნარდობა და სესხის ხელშეკრულებები", 6, 2},
			{"სისხლის სამართლის ზოგადი ნაწილი", "დანაშაულის შემადგენლობა, მართლწინააღმდეგობა, ბრალი და სასჯელი", 6, 2},
			{"სისხლის სამართლის კერძო ნაწილი", "დანაშაულები სიცოცხლის, საკუთრებისა და სახელმწიფოს წინააღმდეგ", 6, 2},
			{"ზოგადი ადმინისტრაციული სამართალი", "ადმინისტრაციული ორგანოები, ინდივიდუალური აქტები და ადმინისტრაციული წარმოება", 6, 2},
			{"საერთაშორისო საჯარო სამართალი", "საერთაშორისო ხელშეკრულებები, სახელმწიფოთა პასუხისმგებლობა და გაერო", 5, 2},
			{"სამოქალაქო საპროცესო სამართალი", "სარჩელი, მტკიცებულებები, სასამართლო განხილვა და გასაჩივრება", 6, 2},
			{"სისხლის საპროცესო სამართალი", "გამოძიება, აღკვეთის ღონისძიებები, ნაფიც მსაჯულთა სასამართლო", 6, 2},
			{"ადმინისტრაციული საპროცესო სამართალი", "ადმინისტრაციული დავების განხილვა სასამართლოში", 5, 2},
			{"შრომის სამართალი", "შრომითი ხელშეკრულება, სამუშაო დრო, შვებულება და შრომითი დავები", 5, 2},
			{"საკორპორაციო სამართალი", "მეწარმეთა შესახებ კანონი, შპს, სააქციო საზოგადოება და მართვა", 5, 2},
			{"საოჯახო და მემკვიდრეობითი სამართალი", "ქორწინება, ალიმენტი, კანონით და ანდერძით მემკვიდრეობა", 5, 2},
			{"ადამიანის უფლებათა საერთაშორისო სამართალი", "ევროპული კონვენცია და სტრასბურგის სასამართლოს პრეცედენტები", 5, 2},
			{"ევროკავშირის სამართალი", "ევროკავშირის ინსტიტუტები, შიდა ბაზარი და ასოცირების შეთანხმება", 5, 2},
			{"საბანკო და საფინანსო სამართალი", "საფინანსო სექტორის რეგულირება და ეროვნული ბანკის ფუნქციები", 4, 2},
			{"სამართლებრივი წერა და კვლევა", "იურიდიული დოკუმენტების შედგენა და სასამართლო პრაქტიკის ანალიზი", 4, 2},

			// --- Faculty 3: Social & Political Sciences (20 subjects) ---
			{"შესავალი პოლიტიკის მეცნიერებაში", "სახელმწიფო, ხელისუფლება, პოლიტიკური იდეოლოგიები და რეჟიმები", 5, 3},
			{"საერთაშორისო ურთიერთობების თეორია", "რეალიზმი, ლიბერალიზმი, კონსტრუქტივიზმი და გლობალური პოლიტიკა", 6, 3},
			{"შესავალი სოციოლოგიაში", "სოციალური სტრუქტურა, კულტურა, სოციალიზაცია და სოციალური უთანასწორობა", 5, 3},
			{"ზოგადი ფსიქოლოგიის საფუძვლები", "აღქმა, მეხსიერება, აზროვნება, ემოციები და პიროვნების თეორიები", 6, 3},
			{"სოციალური ფსიქოლოგია", "სოციალური გავლენა, სტერეოტიპები, ჯგუფური ქცევა და ატიტუდები", 5, 3},
			{"კვლევის თვისებრივი და რაოდენობრივი მეთოდები", "სოციოლოგიური გამოკითხვა, ინტერვიუ, ფოკუს ჯგუფები და მონაცემთა ანალიზი", 6, 3},
			{"შედარებითი პოლიტიკა", "დემოკრატიული და ავტორიტარული სისტემების შედარებითი ანალიზი", 5, 3},
			{"საერთაშორისო ორგანიზაციები", "გაერო, ნატო, ევროკავშირი, ეუთო და რეგიონული ორგანიზაციები", 5, 3},
			{"დიპლომატიის ისტორია და პრაქტიკა", "დიპლომატიური პროტოკოლი, მოლაპარაკებები და საგარეო პოლიტიკის ანალიზი", 5, 3},
			{"ჟურნალისტიკის საფუძვლები", "ახალი ამბების წერა, რეპორტაჟი, ფაქტების გადამოწმება და მედია სტანდარტები", 5, 3},
			{"მასობრივი კომუნიკაციის თეორიები", "მედიის გავლენა საზოგადოებაზე, დღის წესრიგის შექმნა და პროპაგანდა", 5, 3},
			{"საზოგადოებასთან ურთიერთობა (PR)", "სტრატეგიული კომუნიკაცია, კრიზისული PR და ბრენდის იმიჯი", 5, 3},
			{"კონფლიქტების ანალიზი და მართვა", "კონფლიქტების წარმოშობის მიზეზები, მედიაცია და მშვიდობის მშენებლობა", 5, 3},
			{"საჯარო პოლიტიკის ანალიზი", "პოლიტიკის ფორმულირება, იმპლემენტაცია და შეფასება", 5, 3},
			{"გლობალიზაცია და საერთაშორისო უსაფრთხოება", "თანამედროვე უსაფრთხოების გამოწვევები, კიბერ და ენერგო უსაფრთხოება", 5, 3},
			{"პოლიტიკური ფილოსოფია", "პლატონიდან თანამედროვე მოაზროვნეებამდე: სამართლიანობა და თავისუფლება", 5, 3},
			{"განვითარების ფსიქოლოგია", "ადამიანის კოგნიტური და ემოციური განვითარება ბავშვობიდან ხანდაზმულობამდე", 5, 3},
			{"მედია სამართალი და ეთიკა", "გამოხატვის თავისუფლება, ცენზურა, ჟურნალისტური ეთიკა და კონფიდენციალურობა", 4, 3},
			{"კავკასიის გეოპოლიტიკა", "სამხრეთ კავკასიის რეგიონული უსაფრთხოება, ენერგეტიკული პროექტები და პოლიტიკა", 5, 3},
			{"აკადემიური წერა", "სამეცნიერო ნაშრომის სტრუქტურა, არგუმენტაცია და ციტირების სტანდარტები", 4, 3},
		}

		var subjects []models.Subject
		for _, s := range rawSubjects {
			subj := models.Subject{
				Name:        s.name,
				Description: s.description,
				Credit:      s.credit,
				FacultyID:   faculties[s.facultyIdx].ID,
			}
			if err := tx.FirstOrCreate(&subj, models.Subject{Name: subj.Name}).Error; err != nil {
				return err
			}
			subjects = append(subjects, subj)
		}

		// ======== 4. Seed Subject Prerequisites ========
		prerequisites := []models.SubjectPrerequisite{
			// --- Complex Test Setup for Subject 5 (Backend Go) ---
			// Bucket 1: Requires BOTH Databases (Subject 3) AND Frontend (Subject 4)
			{SubjectId: subjects[5].ID, SubjectPrerequisiteId: subjects[3].ID, BucketId: 1},
			{SubjectId: subjects[5].ID, SubjectPrerequisiteId: subjects[4].ID, BucketId: 1},
			// Bucket 2: Requires ONLY OOP (Subject 1)
			{SubjectId: subjects[5].ID, SubjectPrerequisiteId: subjects[1].ID, BucketId: 2},

			// --- Keep other basic chains for broader testing ---
			{SubjectId: subjects[1].ID, SubjectPrerequisiteId: subjects[0].ID, BucketId: 1},   // OOP -> Intro to Prog
			{SubjectId: subjects[2].ID, SubjectPrerequisiteId: subjects[1].ID, BucketId: 1},   // Data Structures -> OOP
			{SubjectId: subjects[11].ID, SubjectPrerequisiteId: subjects[10].ID, BucketId: 1}, // ML -> AI
			{SubjectId: subjects[18].ID, SubjectPrerequisiteId: subjects[6].ID, BucketId: 1},  // Distributed Systems -> Networks

			// Business chains
			{SubjectId: subjects[21].ID, SubjectPrerequisiteId: subjects[20].ID, BucketId: 1}, // Macro -> Micro
			{SubjectId: subjects[23].ID, SubjectPrerequisiteId: subjects[22].ID, BucketId: 1}, // Managerial Acc -> Financial Acc
			{SubjectId: subjects[26].ID, SubjectPrerequisiteId: subjects[22].ID, BucketId: 1}, // Corp Finance -> Financial Acc
			{SubjectId: subjects[30].ID, SubjectPrerequisiteId: subjects[24].ID, BucketId: 1}, // Strategic Mgmt -> Intro Mgmt

			// Law chains
			{SubjectId: subjects[43].ID, SubjectPrerequisiteId: subjects[42].ID, BucketId: 1}, // Property Law -> Civil General
			{SubjectId: subjects[44].ID, SubjectPrerequisiteId: subjects[42].ID, BucketId: 1}, // Obligations -> Civil General
			{SubjectId: subjects[45].ID, SubjectPrerequisiteId: subjects[44].ID, BucketId: 1}, // Contract Law -> Obligations
			{SubjectId: subjects[47].ID, SubjectPrerequisiteId: subjects[46].ID, BucketId: 1}, // Criminal Special -> Criminal General

			// Social Sciences chains
			{SubjectId: subjects[64].ID, SubjectPrerequisiteId: subjects[63].ID, BucketId: 1}, // Social Psych -> General Psych
			{SubjectId: subjects[66].ID, SubjectPrerequisiteId: subjects[60].ID, BucketId: 1}, // Comparative Politics -> Intro PolSci
		}
		for i := range prerequisites {
			if err := tx.FirstOrCreate(&prerequisites[i], models.SubjectPrerequisite{
				SubjectId:             prerequisites[i].SubjectId,
				SubjectPrerequisiteId: prerequisites[i].SubjectPrerequisiteId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 5. Seed Choosable Subjects (All 80 Subjects) ========
		for i, subj := range subjects {
			subjType := models.SubjectProfessional
			if (i+1)%5 == 0 {
				subjType = models.SubjectFree
			}

			idxInFaculty := i % 20
			band := idxInFaculty / 5

			fromSem := models.Semester(band*2 + 1)
			toSem := fromSem + 2

			if band == 3 {
				fromSem = 6
				toSem = 8
			}

			cs := models.ChoosableSubjects{
				FacultyID:    subj.FacultyID,
				SubjectID:    subj.ID,
				SubjectType:  subjType,
				FromSemester: fromSem,
				ToSemester:   toSem,
			}
			if err := tx.FirstOrCreate(&cs, models.ChoosableSubjects{
				FacultyID: cs.FacultyID,
				SubjectID: cs.SubjectID,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 6. Seed Lecturers (8 Lecturers - 2 per Faculty) ========
		lecturers := []models.Lecturer{
			{FirstName: "ლევან", LastName: "ქავთარაძე", Description: "პროგრამული უზრუნველყოფის ინჟინერიის პროფესორი"},
			{FirstName: "ნინო", LastName: "დოლიძე", Description: "მონაცემთა ბაზებისა და ხელოვნური ინტელექტის ლექტორი"},
			{FirstName: "გიორგი", LastName: "მაისურაძე", Description: "ბიზნესის ადმინისტრირებისა და ეკონომიკის დოქტორი"},
			{FirstName: "თამარ", LastName: "ჩხეიძე", Description: "კორპორაციული ფინანსებისა და მარკეტინგის პროფესორი"},
			{FirstName: "დავით", LastName: "კვარაცხელია", Description: "სამოქალაქო და საკორპორაციო სამართლის ასოცირებული პროფესორი"},
			{FirstName: "მარიამ", LastName: "გელაშვილი", Description: "კონსტიტუციური და საერთაშორისო სამართლის დოქტორი"},
			{FirstName: "ალექსანდრე", LastName: "წერეთელი", Description: "საერთაშორისო ურთიერთობებისა და პოლიტიკის მეცნიერების პროფესორი"},
			{FirstName: "ეკატერინე", LastName: "ლომიძე", Description: "სოციალური ფსიქოლოგიისა და კვლევის მეთოდების ლექტორი"},
		}
		for i := range lecturers {
			if err := tx.FirstOrCreate(&lecturers[i], models.Lecturer{
				FirstName: lecturers[i].FirstName,
				LastName:  lecturers[i].LastName,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 7. Seed Subject Groups (Lecture + Seminar for every Subject) ========
		// Note: Subject index `i` maps to:
		//   subjectGroups[i*2]     -> GroupNo 1 (Lecture)
		//   subjectGroups[i*2 + 1] -> GroupNo 2 (Seminar)
		var subjectGroups []models.SubjectGroup
		timeSlots := []struct {
			from string
			to   string
		}{
			{"09:00", "11:00"},
			{"11:00", "13:00"},
			{"13:00", "15:00"},
			{"15:00", "17:00"},
			{"17:00", "19:00"},
		}

		for i, subj := range subjects {
			facultyIdx := i / 20
			lecturer1 := lecturers[facultyIdx*2]
			lecturer2 := lecturers[facultyIdx*2+1]
			slot1 := timeSlots[i%len(timeSlots)]
			slot2 := timeSlots[(i+1)%len(timeSlots)]
			day1 := int8((i % 5) + 1)
			day2 := int8(((i + 2) % 5) + 1)

			groups := []models.SubjectGroup{
				{
					GroupNo:     1,
					SubjectId:   subj.ID,
					Day:         day1,
					From:        slot1.from,
					To:          slot1.to,
					MaxStudent:  35,
					LectureType: models.Lecture,
					LecturerId:  lecturer1.ID,
				},
				{
					GroupNo:     2,
					SubjectId:   subj.ID,
					Day:         day2,
					From:        slot2.from,
					To:          slot2.to,
					MaxStudent:  25,
					LectureType: models.Seminar,
					LecturerId:  lecturer2.ID,
				},
			}

			for j := range groups {
				if err := tx.FirstOrCreate(&groups[j], models.SubjectGroup{
					SubjectId: groups[j].SubjectId,
					GroupNo:   groups[j].GroupNo,
				}).Error; err != nil {
					return err
				}
				subjectGroups = append(subjectGroups, groups[j])
			}
		}

		// ======== 8. Seed Initial Student Priority List ========
		priorities := []models.StudentPriorityList{
			{Priority: 1, SubjectGroupId: subjectGroups[0].ID, StudentId: users[1].ID},  // Subj 0 (CS)
			{Priority: 2, SubjectGroupId: subjectGroups[2].ID, StudentId: users[1].ID},  // Subj 1 (CS)
			{Priority: 1, SubjectGroupId: subjectGroups[80].ID, StudentId: users[2].ID}, // Subj 40 (Law)
			{Priority: 2, SubjectGroupId: subjectGroups[82].ID, StudentId: users[2].ID}, // Subj 41 (Law)
			{Priority: 1, SubjectGroupId: subjectGroups[40].ID, StudentId: users[3].ID}, // Subj 20 (Business)
			{Priority: 1, SubjectGroupId: subjectGroups[60].ID, StudentId: users[4].ID}, // Subj 30 (Business)
		}
		for i := range priorities {
			if err := tx.FirstOrCreate(&priorities[i], models.StudentPriorityList{
				StudentId:      priorities[i].StudentId,
				SubjectGroupId: priorities[i].SubjectGroupId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 9. Seed Enrollments ========
		enrollments := []models.EnrollmentTable{
			// --- Historical Passed/Failed records (Should NOT count toward active group capacity) ---
			{Status: models.StatusPassed, SubjectGroupID: subjectGroups[0].ID, Year: 2025, Semester: 1, StudentID: users[1].ID},  // g.beridze passed Subj 0
			{Status: models.StatusPassed, SubjectGroupID: subjectGroups[2].ID, Year: 2025, Semester: 1, StudentID: users[1].ID},  // g.beridze passed Subj 1 (Triggers Bucket 2 for Subj 5)
			{Status: models.StatusPassed, SubjectGroupID: subjectGroups[40].ID, Year: 2025, Semester: 2, StudentID: users[3].ID}, // d.maisuradze passed Bus Subj 20
			{Status: models.StatusFailed, SubjectGroupID: subjectGroups[6].ID, Year: 2025, Semester: 1, StudentID: users[12].ID}, // e.kvaratskhelia failed CS Subj 3 Group 1

			// --- Active CS Enrollments (Tested by cs.manager@university.edu) ---
			// CS Subject 0 (Python): Group 1 gets 3 active students, Group 2 gets 1 active student
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[0].ID, Year: 2026, Semester: 1, StudentID: users[9].ID},
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[0].ID, Year: 2026, Semester: 1, StudentID: users[10].ID},
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[0].ID, Year: 2026, Semester: 1, StudentID: users[11].ID},
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[1].ID, Year: 2026, Semester: 1, StudentID: users[12].ID},

			// CS Subject 1 (Java OOP): Group 1 gets 2 active students, Group 2 has 0 active students
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[2].ID, Year: 2026, Semester: 2, StudentID: users[9].ID},
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[2].ID, Year: 2026, Semester: 2, StudentID: users[10].ID},

			// CS Subject 3 (SQL Databases): Group 1 gets 1 active student (plus 1 failed above -> should return count 1)
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[6].ID, Year: 2026, Semester: 2, StudentID: users[11].ID},

			// --- Active Business Enrollments (Tested by bus.manager@university.edu) ---
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[40].ID, Year: 2026, Semester: 1, StudentID: users[3].ID}, // Bus Subj 20 Group 1 (passed in 2025, retaking/active or new group)
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[41].ID, Year: 2026, Semester: 1, StudentID: users[4].ID}, // Bus Subj 20 Group 2
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[42].ID, Year: 2026, Semester: 2, StudentID: users[4].ID}, // Bus Subj 21 Group 1

			// --- Active Law Enrollment ---
			{Status: models.StatusActive, SubjectGroupID: subjectGroups[80].ID, Year: 2026, Semester: 1, StudentID: users[2].ID}, // Law Subj 40 Group 1
		}
		for i := range enrollments {
			if err := tx.FirstOrCreate(&enrollments[i], models.EnrollmentTable{
				StudentID:      enrollments[i].StudentID,
				SubjectGroupID: enrollments[i].SubjectGroupID,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 10. Seed Lecturer Ratings ========
		for i, lec := range lecturers {
			rating := models.LecturerRating{
				LecturerId: lec.ID,
				Rating:     int16(4 + (i % 2)),
			}
			if err := tx.FirstOrCreate(&rating, models.LecturerRating{
				LecturerId: rating.LecturerId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 11. Seed Subject Recommendations ========
		recommendations := []models.SubjectRecommendation{
			{SubjectId: subjects[0].ID, UserId: users[1].ID, VoteType: models.Like},
			{SubjectId: subjects[1].ID, UserId: users[1].ID, VoteType: models.Like},
			{SubjectId: subjects[20].ID, UserId: users[3].ID, VoteType: models.Like},
			{SubjectId: subjects[40].ID, UserId: users[2].ID, VoteType: models.Like},
		}
		for i := range recommendations {
			if err := tx.FirstOrCreate(&recommendations[i], models.SubjectRecommendation{
				SubjectId: recommendations[i].SubjectId,
				UserId:    recommendations[i].UserId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 12. Seed Subject Difficulties (All 80 Subjects) ========
		for i, subj := range subjects {
			diff := models.SubjectDifficulty{
				SubjectId:       subj.ID,
				DifficultyScale: int16(2 + (i % 4)),
			}
			if err := tx.FirstOrCreate(&diff, models.SubjectDifficulty{
				SubjectId: diff.SubjectId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 13. Seed Registration Windows ========
		activeReg := models.ActiveRegistration{
			From: time.Now(),
			To:   time.Now().Add(14 * 24 * time.Hour),
		}
		if err := tx.Create(&activeReg).Error; err != nil {
			return err
		}

		inactiveReg := models.InactiveRegistration{
			From: time.Now().Add(-30 * 24 * time.Hour),
			To:   time.Now().Add(-15 * 24 * time.Hour),
		}
		if err := tx.Create(&inactiveReg).Error; err != nil {
			return err
		}

		// ======== 14. Seed Student Subject Allowances ========
		for _, user := range users {
			if user.Role != models.RoleStudent {
				continue
			}

			for sem := models.Semester(1); sem <= 8; sem++ {
				profAllowance := int(1 + (sem % 2))
				freeAllowance := int(2 - (sem % 2))

				allowance := models.StudentSubjectAllowance{
					StudentId:                    user.ID,
					ProfessionalSubjectAllowance: profAllowance,
					FreeSubjectAllowance:         freeAllowance,
					Semester:                     sem,
				}

				if err := tx.FirstOrCreate(&allowance, models.StudentSubjectAllowance{
					StudentId: allowance.StudentId,
					Semester:  allowance.Semester,
				}).Error; err != nil {
					return err
				}
			}
		}

		log.Println("[INFO] Database seeded successfully with 4 faculties, 80 subjects, 3 managers, 9 students, and active/historical enrollment scenarios!")
		return nil
	})
}