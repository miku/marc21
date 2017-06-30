// Package marc21 implements support for the MARC21 bibliographic format.
//
// Copyright (C) 2011 William Waites
// Copyright (C) 2012 Dan Scott <dan@coffeecode.net>
// Copyright (C) 2017 Martin Czygan <martin.czygan@uni-leipzig.de>
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of
// the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License and the GNU General Public License along with this program
// (the files COPYING and GPL3 respectively).  If not, see
// <http://www.gnu.org/licenses/>.
//
// Package marc21 reads and writes MARC21 bibliographic catalogue records.
//
// Usage is straightforward. For example,
//
//     marcfile, err := os.Open("somedata.mrc")
//     record, err := marc21.ReadRecord(marcfile)
//     err = record.WriteTo(os.Stdout)
package marc21
