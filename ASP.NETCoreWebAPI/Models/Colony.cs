namespace ASP.NETCoreWebAPI.Models
{
    public class Colony
    {
        public int Id { get; set; } // Primary Key
        public int AccountLevel { get; set; } // Attribute
        public ColonyAsset[]? Assets { get; set; } // Multi-Valued Foreign Key (ColonyAsset) - (Should this be 'nullable (?)' or 'required'?)
        public required ColonyCode colonyCode { get; set; } // Foreign Key (ColonyCode) - (Should this be 'nullable (?)' or 'required'?)
        public int Owner { get; set; } // Foreign Key (Player)
        public Location[]? Locations { get; set; } // Multi-Valued Foreign Key (Location) - (Should this be 'nullable (?)' or 'required'?)
        public DateTime[]? LastestVisit { get; set; } // Multi-Valued Attribute
    }
}
