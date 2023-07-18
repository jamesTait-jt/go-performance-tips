package loops

type bigObj struct {
	id int
	lotsOfStuff [10 * 1024]byte
}

type otherObj struct {
	id int
}

func forI(sl []bigObj) otherObj {
	oo := otherObj{}
	for i := 0 ; i < len(sl) ; i++ {
		if oo.id == 5 {
			oo.id = sl[i].id
		}
	}
	return oo 
}

func forRange(sl[]bigObj) otherObj {
	oo := otherObj{}
	for _, elem := range sl {
		if oo.id == 5 {
			oo.id = elem.id
		}

	}
	return oo
}

func forINoIf(sl []bigObj) otherObj {
	oo := otherObj{}
	for i := 0 ; i < len(sl) ; i++ {
		oo.id = sl[i].id
	}
	return oo 
}

func forRangeNoIf(sl[]bigObj) otherObj {
	oo := otherObj{}
	for _, elem := range sl {
		oo.id = elem.id
	}
	return oo
}