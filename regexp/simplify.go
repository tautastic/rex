package regexp

func simplifyRegexp(re *Regexp) *Regexp {
	if re.Op != OpConcat {
		return &Regexp{Op: OpConcat, Sub: []*Regexp{re, {Op: OpAccept}}}
	}
	re.Sub = append(re.Sub, &Regexp{Op: OpAccept})
	return re
}
